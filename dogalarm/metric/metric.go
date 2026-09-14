package metric

import "github.com/prometheus/client_golang/prometheus"

var (
	DropAlarmCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "drop_alarm_total",
		},
		[]string{"alarm_app", "alarm_host", "notice_type"},
	)
	LiveProbeGuage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "live_probe",
		},
		[]string{"alarm_app", "alarm_host"},
	)
)

type LiveStatus int 

const (
	Living LiveStatus = 1
	ShutDown LiveStatus = 2
)

func All() []prometheus.Collector {
	return []prometheus.Collector{
		DropAlarmCounter,
		LiveProbeGuage,
	}
}