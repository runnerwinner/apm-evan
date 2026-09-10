package dogapm

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type redisHook struct {
	
}

const (
	redisTracerName = "dogapm/redis"
)

func (r *redisHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (r *redisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	tracer := otel.Tracer(redisTracerName)
	return func(ctx context.Context, cmd redis.Cmder) error {
		if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
			fmt.Printf("[redis hook] traceID=%s cmd=%s\n", span.SpanContext().TraceID().String(), cmd.String())
		} else {
			fmt.Printf("[redis hook] no active span, cmd=%s\n", cmd.String())
		}
		ctx, span := tracer.Start(ctx, "redisProcessCmd")
		span.SetAttributes(attribute.String("cmd", truncate(cmd.String())))
		defer span.End()
		err := next(ctx, cmd)
		if err != nil && err != redis.Nil {
			span.SetAttributes(attribute.Bool("error", true))
			span.RecordError(err, trace.WithStackTrace(true))
		}
		return err
	}
}


func (r *redisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	tracer := otel.Tracer(redisTracerName)
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
			fmt.Printf("[redis hook pipeline] traceID=%s cmds=%v\n", span.SpanContext().TraceID().String(), cmds)
		} else {
			fmt.Printf("[redis hook pipeline] no active span, cmds=%v\n", cmds)
		}
		ctx, span := tracer.Start(ctx, "redisProcessPipeline")
		span.SetAttributes(attribute.String("cmd", truncate(fmt.Sprintf("%v",cmds))))
		defer span.End()
		err := next(ctx, cmds)
		if err != nil && err != redis.Nil {
			span.SetAttributes(attribute.Bool("error", true))
			span.RecordError(err, trace.WithStackTrace(true))
		}
		return err
	}
}