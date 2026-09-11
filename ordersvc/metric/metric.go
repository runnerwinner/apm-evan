package metric

import "github.com/prometheus/client_golang/prometheus"
var (
	OrderSuccessCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_success_total",
		}, []string{"sku_id"},
	)
	// OrderFailedCounter = prometheus.NewCounterVec(
	// 	prometheus.CounterOpts{
	// 		Name: "order_failed_total",
	// 	}, []string{"type", "method","peer","peer_host"},
	// )
	// OrderHandleHistogram = prometheus.NewHistogramVec(
	// 	prometheus.HistogramOpts{
	// 		Name: "order_handle_seconds",
	// 	}, []string{"type", "method", "status","peer","peer_host"},
	// )
)

func All() []prometheus.Collector {
	return []prometheus.Collector{
		OrderSuccessCounter,
		// OrderFailedCounter,
		// OrderHandleHistogram,
	}
}