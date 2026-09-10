package main

import (
	"dogapm"
	"os"
	"protos"
	"skusvc/grpc"
	"time"
)

func main() {
	_ = os.Setenv("OTEL_SERVICE_NAME", "skusvc")

	//初始化db, http server, grpcclient

	dogapm.Infra.Init(
		dogapm.InfraDbOption("root:password@tcp(localhost:3307)/skusvc"),
		dogapm.InfraEnableApm("127.0.0.1:54317", 15*time.Second),
	)

	grpcserver := dogapm.NewGrpcServer(":8001")
	protos.RegisterSkuServiceServer(grpcserver, &grpc.SkuServer{})

	// httpserver := dogapm.NewHttpServer(":8081")
	// httpserver.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	_, _ = w.Write([]byte("OK"))
	// })
	// httpserver.HandleFunc("/sku/add", api.Sku.Add)

	dogapm.EndPoint.Start()

}