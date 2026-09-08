package main

import (
	"dogapm"
	"protos"
	"skusvc/grpc"
)

func main() {

	//初始化db, http server, grpcclient

	dogapm.Infra.Init(
		dogapm.InfraDbOption("root:password@tcp(localhost:3307)/skusvc"),
		// dogapm.InfraRdbOption("localhost:6380"),
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