package main

import (
	"dogapm"
	"os"
	"protos"
	"time"
	"usrsvc/grpc"
)

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	if os.Getenv("OTEL_SERVICE_NAME") == "" {
		_ = os.Setenv("OTEL_SERVICE_NAME", "usrsvc")
	}

	//初始化db, http server, grpcclient
	dbDSN := envOrDefault("DB_DSN", "root:password@tcp(localhost:3307)/usrsvc")
	redisAddr := envOrDefault("REDIS_ADDR", "localhost:6380")
	otelAddr := envOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "127.0.0.1:54317")

	dogapm.Infra.Init(
		dogapm.InfraDbOption(dbDSN),
		dogapm.InfraRdbOption(redisAddr),
		dogapm.InfraEnableApm(otelAddr, 15*time.Second),
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