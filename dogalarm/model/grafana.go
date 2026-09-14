package model


const (
	StatusFiring = "firing"
	StatusAlert = "alert"
	StatusResolved = "resolved"
)

type AlertModel struct {
	Receiver string `json:"receiver"`
	Status string `json:"status"`
	Alerts []Alerts `json:"alerts"`
	GroupLabels GroupLabels `json:"groupLabels"`
	CommonLabels map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL string `json:"externalURL"`
	Version string `json:"version"`
	GroupKey string `json:"groupKey"`
	TruncatedAlerts int `json:"truncatedAlerts"`
	OrgID int `json:"orgId"`
	Title string `json:"title"`
	State string `json:"state"`
	Message string `json:"message"`
}



type GroupLabels struct {
	AlertName string `json:"alertname"`
	GrafanaFolder string `json:"grafana_folder"`
}

type Alerts struct {
	Status string `json:"status"`
	Labels map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt string `json:"startsAt"`
	EndsAt string `json:"endsAt"`
	GeneratorURL string `json:"generatorURL"`
	Fingerprint string `json:"fingerprint"`
	SilenceURL string `json:"silenceURL"`
	DashboardURL string `json:"dashboardURL"`
	PanelURL string `json:"panelURL"`
	Values map[string]any `json:"values"`
	ValueString string `json:"valueString"`
}