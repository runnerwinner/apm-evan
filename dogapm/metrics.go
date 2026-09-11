package dogapm

import (
	"dogapm/internal"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

const (
	TypeHttp = "http"
	TypeGrpc = "grpc"
	TypeMysql = "mysql"
	TypeRedis = "redis"
)

type customMetricRegistry struct {
	*prometheus.Registry
	customLabels []*io_prometheus_client.LabelPair
}

func newCustomMetricRegistry(customLabels map[string]string) *customMetricRegistry {
	c := &customMetricRegistry{
		Registry:     prometheus.NewRegistry(),		
	}
	for k, v := range customLabels {
		kCp := k
		vCp := v
		c.customLabels = append(c.customLabels, &io_prometheus_client.LabelPair{
			Name:  &kCp,
			Value: &vCp,
		})
	}
	return c
}

func (c *customMetricRegistry) Gather() ([]*io_prometheus_client.MetricFamily, error) {
	metricFamilies, err := c.Registry.Gather()
	if err != nil {
		return nil, err
	}
	for _, mf := range metricFamilies {
		for _, m := range mf.Metric {
			m.Label = append(m.Label, c.customLabels...)
		}
	}
	return metricFamilies, nil
}

var (
	MetricsReg = newCustomMetricRegistry(map[string]string{
		"host": internal.BuildInfo.Hostname(),
		"app":  internal.BuildInfo.AppName(),
	})
)

var (

	serverHandleHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "server_handle_seconds",
		}, []string{"type", "method", "status"},
	)

	serverHandleCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "server_handle_total",
		}, []string{"type", "method"},
	)

	clientHandleHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "client_handle_seconds",
		}, []string{"type", "method", "server"},
	)

	clientHandleCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "client_handle_total", 
		}, []string{"type", "method", "server"},
	)

	// 对 mysql 或 redis 等第三方库的调用次数进行统计
	// method : mysql - query delete, redis - get set  and so on
	// name ： 表名 
	// server : 标识第三方组件的一些信息， 比如 mysql - host:port, redis - host:port
	libraryCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lib_handle_total",
		}, []string{"type", "method", "name", "server"},
	)
)

func init() {
	MetricsReg.MustRegister(serverHandleHistogram,serverHandleCounter,clientHandleHistogram,clientHandleCounter,libraryCounter)
}
