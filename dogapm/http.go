package dogapm

import (
	"context"
	"errors"
	"net"
	"net/http"
)


type HttpServer struct {
	mux *http.ServeMux
	*http.Server
}


func NewHttpServer(addr string) *HttpServer {
	mux := http.NewServeMux()
	server := &HttpServer{
		mux: mux,
		Server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
	globalClosers = append(globalClosers, server)
	globalStarters = append(globalStarters, server)
	return server
}


func (s *HttpServer) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *HttpServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
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