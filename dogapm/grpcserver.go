package dogapm

import (
	"context"
	"net"

	"google.golang.org/grpc"
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

func unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Pre-processing logic before handling the RPC
		resp, err := handler(ctx, req)
		// Post-processing logic after handling the RPC
		return resp, err
	}
}