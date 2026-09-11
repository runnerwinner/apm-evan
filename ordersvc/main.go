package main

import (
	"dogapm"
	"net/http"
	"ordersvc/api"
	"ordersvc/grpcclient"
	"ordersvc/metric"
	"os"
	"protos"
	"time"
)

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	if os.Getenv("OTEL_SERVICE_NAME") == "" {
		_ = os.Setenv("OTEL_SERVICE_NAME", "ordersvc")
	}

	//初始化db, http server, grpcclient
	dbDSN := envOrDefault("DB_DSN", "root:password@tcp(localhost:3307)/ordersvc")
	otelAddr := envOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "127.0.0.1:54317")
	skuAddr := envOrDefault("SKU_GRPC_ADDR", "localhost:8001")
	userAddr := envOrDefault("USER_GRPC_ADDR", "localhost:8002")

	dogapm.Infra.Init(
		dogapm.InfraDbOption(dbDSN),
		dogapm.InfraEnableApm(otelAddr, 15*time.Second),
		dogapm.MetricReg(metric.All()...),
	)

	// TODO: grpcclient初始化
	skuconn, err := dogapm.NewGrpcClient(skuAddr,"skusvc")
	if err != nil {
		panic(err)
	}
	grpcclient.SkuClient = protos.NewSkuServiceClient(skuconn)

	userconn, err := dogapm.NewGrpcClient(userAddr,"usrsvc")
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
