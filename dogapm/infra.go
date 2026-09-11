package dogapm

import (
	"context"
	"database/sql"
	"dogapm/internal"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type infra struct {
	Db  *sql.DB
	Rdb *redis.Client
}

var Infra = &infra{}

type InfraOption func(i *infra)

func MetricReg(collectors ...prometheus.Collector) InfraOption{
	return func(i *infra) {
		MetricsReg.MustRegister(collectors...)
	}
}

func InfraDbOption(connectUrl string) InfraOption {
	return func(i *infra) {
		db, err := sql.Open(wrappedMySQLDriver(connectUrl), connectUrl)
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		i.Db = db
		err = db.Ping()
		if err != nil {
			panic(err)
		}
	}
}

func InfraRdbOption(connectUrl string) InfraOption {
	return func(i *infra) {
		rdb := redis.NewClient(&redis.Options{
			Addr: connectUrl,
			DB:   0,
		})
		rdb.AddHook(&redisHook{})
		if err := rdb.Ping(context.TODO()).Err(); err != nil {
			panic(err)
		}
		i.Rdb = rdb
	}
}

func InfraEnableApm(otelEndpoint string, connectTimeout ...time.Duration) InfraOption {
	return func(i *infra) {
		ctx := context.Background()
		res, err := resource.New(ctx, resource.WithAttributes(
			semconv.ServiceName(internal.BuildInfo.AppName()),
		))
		if err != nil {
			panic(fmt.Errorf("create otel resource failed: %w", err))
		}
		timeout := 5 * time.Second
		if len(connectTimeout) > 0 && connectTimeout[0] > 0 {
			timeout = connectTimeout[0]
		}
		grpcTarget := otelEndpoint
		if !strings.Contains(otelEndpoint, "://") {
			grpcTarget = "passthrough:///" + otelEndpoint
		}
		probeCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		conn, err := grpc.NewClient(grpcTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Errorf("create otel gRPC client for %q failed: %w", grpcTarget, err))
		}

		// Set up a trace exporter
		traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			panic(fmt.Errorf("create otel trace exporter for %q failed: %w", grpcTarget, err))
		}

		bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
		tracerProvider := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(bsp),
			sdktrace.WithResource(res),
		)
		otel.SetTracerProvider(tracerProvider)
		probeTracer := tracerProvider.Tracer(internal.BuildInfo.AppName())
		_, probeSpan := probeTracer.Start(ctx, "otel-startup-probe")
		probeSpan.End()
		if err := tracerProvider.ForceFlush(probeCtx); err != nil {
			panic(fmt.Errorf("flush otel startup probe to %q within %s failed: %w", grpcTarget, timeout, err))
		}

		// Set the global propagator to tracecontext (the default is no-op).
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		globalClosers = append(globalClosers, &traceProviderComponent{provider: tracerProvider})

	}
}

type traceProviderComponent struct {
	provider *sdktrace.TracerProvider
}

func (t *traceProviderComponent) Close() error {
	if t.provider != nil {
		if err := t.provider.Shutdown(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

func (i *infra) Init(options ...InfraOption) {
	for _, option := range options {
		option(i)
	}
	Tracer = otel.Tracer(internal.BuildInfo.AppName())
}
