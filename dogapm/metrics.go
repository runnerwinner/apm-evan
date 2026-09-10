package dogapm

import (
	"dogapm/internal"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

const (
	TypeHttp = "http"
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
)

func init() {
	MetricsReg.MustRegister(serverHandleHistogram,serverHandleCounter)
}
