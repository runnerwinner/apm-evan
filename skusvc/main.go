package main

import (
	"dogapm"
	"os"
	"protos"
	"skusvc/grpc"
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
		_ = os.Setenv("OTEL_SERVICE_NAME", "skusvc")
	}

	//初始化db, http server, grpcclient
	dbDSN := envOrDefault("DB_DSN", "root:password@tcp(localhost:3307)/skusvc")
	otelAddr := envOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "127.0.0.1:54317")

	dogapm.Infra.Init(
		dogapm.InfraDbOption(dbDSN),
		dogapm.InfraEnableApm(otelAddr, "/logs", 2, 15*time.Second),
	)

	dogapm.NewHttpServer(":8091")

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