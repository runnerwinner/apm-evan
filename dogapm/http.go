package dogapm

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	httpTracerName = "dogapm/http"
)

type HttpServer struct {
	mux *http.ServeMux
	*http.Server
	tracer trace.Tracer
}


func NewHttpServer(addr string) *HttpServer {
	mux := http.NewServeMux()
	server := &HttpServer{
		mux: mux,
		Server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		tracer: otel.Tracer(httpTracerName),
	}
	globalClosers = append(globalClosers, server)
	globalStarters = append(globalStarters, server)
	return server
}


func (s *HttpServer) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, &traceHandler{
		handler: handler,
		tracer:  s.tracer,
	})
}

func (s *HttpServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.Handle(pattern, &traceHandler{
		handler: http.HandlerFunc(handler),
		tracer:  s.tracer,
	})
}


type traceHandler struct {
	handler http.Handler
	tracer trace.Tracer
}

type respWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (w *respWriterWrapper) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (t *traceHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if t.tracer == nil {
		t.handler.ServeHTTP(writer, request)
		return
	}
	ctx := request.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(request.Header))
	ctx, span := t.tracer.Start(ctx, request.URL.Path)
	defer span.End()
	request = request.Clone(ctx)
	start := time.Now()
	respWrapper := &respWriterWrapper{ResponseWriter: writer}
	t.handler.ServeHTTP(respWrapper, request)
	if respWrapper.status == 0 {
		respWrapper.status = http.StatusOK
	}
	end := time.Now()
	span.SetAttributes(
		attribute.KeyValue{
			Key:   "http.status_code",
			Value: attribute.StringValue(strconv.Itoa(respWrapper.status)),
		},
		attribute.KeyValue{
			Key:   "http.duration",
			Value: attribute.IntValue(int(end.Sub(start).Seconds())),
		},
	)
}

func (s *HttpServer) Start() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	go func() {
		if err := s.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// handle error, e.g., log it
			panic(err)
		}
	}()
	return nil
}

func (s *HttpServer) Close() error {
	return s.Shutdown(context.Background())
}