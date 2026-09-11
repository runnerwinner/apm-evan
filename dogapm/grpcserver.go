package dogapm

import (
	"context"
	"net"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	*grpc.Server
	addr string
}

func NewGrpcServer(addr string) *GrpcServer {
	svc := grpc.NewServer(grpc.UnaryInterceptor(unaryServerInterceptor()))
	server := &GrpcServer{
		Server: svc,
		addr:   addr,
	}
	globalClosers=append(globalClosers, server)
	globalStarters=append(globalStarters, server)
	return server
}

func (s *GrpcServer) Close() error {
	s.GracefulStop()
	return nil
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	go func() {
		if err := s.Serve(lis); err != nil {
			// handle error, e.g., log it
			panic(err)
		}
	}()
	return nil
}

const (
	grpcServerTracerName = "dogapm/grpc_server"
)

func unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Pre-processing logic before handling the RPC
		tracer := otel.Tracer(grpcServerTracerName)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		ctx = otel.GetTextMapPropagator().Extract(ctx, &metadataSupplier{metadata: &md})
		ctx, span := tracer.Start(ctx, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
		start := time.Now()
		statusCode := codes.OK
		defer func() {			
			span.SetAttributes(attribute.Float64("grpc.duration", float64(time.Since(start).Seconds())))
			span.End()
			serverHandleHistogram.WithLabelValues(TypeGrpc, info.FullMethod, strconv.Itoa(int(statusCode))).Observe(time.Since(start).Seconds())
		}()
		serverHandleCounter.WithLabelValues(TypeGrpc, info.FullMethod).Inc()
		resp, err := handler(ctx, req)
		if err != nil {
			s, _ := status.FromError(err)
			statusCode = s.Code()
			span.RecordError(err,trace.WithTimestamp(time.Now()), trace.WithStackTrace(true))
			span.SetAttributes(attribute.Bool("error",true))
			span.SetAttributes(attribute.String("grpc.status_code", s.Code().String()))
		}
		
		// Post-processing logic after handling the RPC
		return resp, err
	}
}