package dogapm

import (
	"context"
	"maps"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)
type log struct{}

const traceId = "traceId"

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

var Logger = &log{}

func (l *log) buildFields(action string, kv map[string]any) map[string]any {
	fields := make(map[string]any, len(kv)+1)
	maps.Copy(fields, kv)
	fields["action"] = action
	return fields
}

func (l *log) Info(ctx context.Context, action string, kv map[string]any) {
	logrus.WithFields(l.buildFields(action, kv)).Info()
}

func (l *log) Warn(ctx context.Context, action string, kv map[string]any) {
	logrus.WithFields(l.buildFields(action, kv)).Warn()
}

func (l *log) Debug(ctx context.Context, action string, kv map[string]any) {
	logrus.WithFields(l.buildFields(action, kv)).Debug()
}

func (l *log) Error(ctx context.Context, action string, kv map[string]any, err error) {
	fields := l.buildFields(action, kv)
	if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
		fields[traceId] = span.SpanContext().TraceID().String()
		span.SetAttributes(attribute.Bool("error", true))
		span.RecordError(err, trace.WithStackTrace(true), trace.WithTimestamp(time.Now()))
	}
	logrus.WithFields(fields).WithError(err).Error()
}