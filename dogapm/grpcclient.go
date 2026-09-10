package dogapm

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GrpcClient wraps a *grpc.ClientConn. The embedded ClientConn promotes the
// full gRPC API for direct use, while this type serves as an extension point
// for APM hooks (tracing, metrics, connection lifecycle).
type GrpcClient struct {
	*grpc.ClientConn
}

// NewGrpcClient creates a lazily connected gRPC client to addr.
//
// It uses grpc.NewClient instead of the deprecated grpc.Dial (deprecated
// since gRPC-Go v1.63). The connection is NOT established until the first RPC,
// so an unreachable addr will not fail here - it surfaces as an RPC error.
// Always set a timeout on the RPC context to avoid hanging, or call
// conn.Connect()/WaitForStateChange if a fail-fast startup check is required.
//
// opts are appended after the defaults so callers can override the insecure
// transport credentials used by default for development, e.g. pass
// grpc.WithTransportCredentials(tlsCredentials) in production.
func NewGrpcClient(addr string, opts ...grpc.DialOption) (*GrpcClient, error) {
	dialOpts := append([]grpc.DialOption{
		grpc.WithUnaryInterceptor(unaryInterceptor()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, opts...)

	conn, err := grpc.NewClient(addr, dialOpts...)
	if err != nil {
		return nil, err
	}
	return &GrpcClient{ClientConn: conn}, nil
}

const (
	grpcClientTracerName = "dogapm/grpc_client"
)

// unaryInterceptor is the extension point where APM hooks (trace injection,
// metric recording, error tagging) should be attached on the client side.
func unaryInterceptor() grpc.UnaryClientInterceptor {
	tracer := otel.Tracer(grpcClientTracerName)
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Pre-processing logic before invoking the RPC
		ctx,span := tracer.Start(ctx, method, trace.WithSpanKind(trace.SpanKindClient))
		start := time.Now()
		defer func() {
			span.SetAttributes(attribute.Float64("grpc.duration", float64(time.Now().Sub(start).Seconds())))
			span.End()
			
		}()
		md,ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.MD{}		
		}
		
		otel.GetTextMapPropagator().Inject(ctx, &metadataSupplier{metadata: &md})
		err := invoker(ctx, method, req, reply, cc, opts...)
		// Post-processing logic after invoking the RPC
		if err != nil {
			s, _ := status.FromError(err)
			span.RecordError(err,trace.WithTimestamp(time.Now()), trace.WithStackTrace(true))
			span.SetAttributes(attribute.Bool("error",true))
			span.SetAttributes(attribute.String("grpc.status_code", s.Code().String()))
		}
		
		return err
	}
}