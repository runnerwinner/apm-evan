package main

import (
	"dogapm"
	"os"
	"protos"
	"time"
	"usrsvc/grpc"
)

func main() {
	_ = os.Setenv("OTEL_SERVICE_NAME", "usrsvc")

	//初始化db, http server, grpcclient

	dogapm.Infra.Init(
		dogapm.InfraDbOption("root:password@tcp(localhost:3307)/usrsvc"),
		dogapm.InfraRdbOption("localhost:6380"),
		dogapm.InfraEnableApm("127.0.0.1:54317", 15*time.Second),
	)

	dogapm.NewHttpServer(":8092")

	grpcserver := dogapm.NewGrpcServer(":8002")
	protos.RegisterUserServiceServer(grpcserver, &grpc.UserServer{})

	// httpserver := dogapm.NewHttpServer(":8081")
	// httpserver.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	_, _ = w.Write([]byte("OK"))
	// })
	// httpserver.HandleFunc("/sku/add", api.Sku.Add)

	dogapm.EndPoint.Start()

}