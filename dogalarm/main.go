package main

import (
	"dogalarm/api"
	"dogalarm/metric"
	"dogapm"
	"net/http"
	"os"
	"time"
)

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	otelAddr := envOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "127.0.0.1:54317")
	dbDSN := envOrDefault("DB_DSN", "root:password@tcp(localhost:3307)/dogalarm")
	redisAddr := envOrDefault("REDIS_ADDR", "localhost:6380")
	dogapm.Infra.Init(
		dogapm.InfraEnableApm(otelAddr, "/logs", 2, 15*time.Second),
		dogapm.InfraDbOption(dbDSN),
		dogapm.InfraRdbOption(redisAddr),
		dogapm.MetricReg(metric.All()...),
	)

	httpServer := dogapm.NewHttpServer(":8084")
	httpServer.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	httpServer.HandleFunc("/metric_webhook", api.Alarm.MetricWebHook)
	httpServer.HandleFunc("/log_webhook", api.Alarm.LogWebHook)
	dogapm.EndPoint.Start()
}
