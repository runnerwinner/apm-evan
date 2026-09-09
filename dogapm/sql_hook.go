package dogapm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	mysqlTracerName        = "dogapm/mysql"
	ctxKeyBeginTime        = "begin"
	wrappedMySQLDriverName = "dogapm-mysql"
)

var registerWrappedMySQLDriverOnce sync.Once

func wrappedMySQLDriver() string {
	registerWrappedMySQLDriverOnce.Do(func() {
		sql.Register(wrappedMySQLDriverName, wrap(mysql.MySQLDriver{}))
	})
	return wrappedMySQLDriverName
}

func truncate(query string) string {
	if len(query) > 1024 {
		return query[:1024]
	}
	return query
}

func wrap(d driver.Driver) driver.Driver {

	tracer := otel.Tracer(mysqlTracerName)
	return &Driver{
		Driver: d,
		hooks: Hooks{
			Before: func(ctx context.Context, query string, args ...any) (context.Context, error) {
				ctx = context.WithValue(ctx, ctxKeyBeginTime, time.Now())
				if ctx, span := tracer.Start(ctx, "sqltrace"); span != nil {
					span.SetAttributes(
						attribute.String("sql", truncate(query)),
						attribute.String("param", truncate(fmt.Sprintf("%v", args...))),
					)
					return ctx, nil
				}
				return ctx, nil
			},
			After: func(ctx context.Context, query string, args ...any) (context.Context, error) {
				beginTime := time.Now()
				if v := ctx.Value(ctxKeyBeginTime); v != nil {
					if bt, ok := v.(time.Time); ok {
						beginTime = bt
					}
				}
				now := time.Now()
				span := trace.SpanFromContext(ctx)
				if now.Sub(beginTime).Seconds() > 1 {
					span.SetAttributes(attribute.Bool("slowsql", true))
				}
				span.End()
				return ctx, nil
			},
			OnError: func(ctx context.Context, err error, query string, args ...any) error {
				span := trace.SpanFromContext(ctx)
				if !errors.Is(err, driver.ErrSkip) {
					span.SetAttributes(
						attribute.Bool("error", true),
					)
					span.RecordError(err, trace.WithStackTrace(true))
				} else {
					span.SetAttributes(attribute.Bool("drop", true))
				}
				span.End()
				return err
			},
		},
	}
}
