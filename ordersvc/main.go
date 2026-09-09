package main

import (
	"dogapm"
	"net/http"
	"ordersvc/api"
	"ordersvc/grpcclient"
	"protos"
	"time"
)

func main() {

	//初始化db, http server, grpcclient

	dogapm.Infra.Init(
		dogapm.InfraDbOption("root:password@tcp(localhost:3307)/ordersvc"),
		dogapm.InfraEnableApm("127.0.0.1:54317", 15*time.Second),
	)

	// TODO: grpcclient初始化
	skuconn, err := dogapm.NewGrpcClient("localhost:8001")
	if err != nil {
		panic(err)
	}
	grpcclient.SkuClient = protos.NewSkuServiceClient(skuconn)
	userconn, err := dogapm.NewGrpcClient("localhost:8002")
	if err != nil {
		panic(err)
	}
	grpcclient.UserClient = protos.NewUserServiceClient(userconn)

	httpserver := dogapm.NewHttpServer(":8080")
	httpserver.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	httpserver.HandleFunc("/order/add", api.Order.Add)

	dogapm.EndPoint.Start()
}
