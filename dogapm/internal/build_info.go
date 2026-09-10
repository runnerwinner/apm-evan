package internal

import (
	"os"
	"path/filepath"
)
var(
	hostname string
	appName string
)

func init() {
	hostname,_ = os.Hostname()
	appName = filepath.Base(os.Args[0])
	if serviceName := os.Getenv("OTEL_SERVICE_NAME"); serviceName != "" {
		appName = serviceName
	}
}

type buildInfo struct {

}

var BuildInfo = &buildInfo{}

func (b *buildInfo) Hostname() string {
	return hostname
}

func (b *buildInfo) AppName() string {
	if serviceName := os.Getenv("OTEL_SERVICE_NAME"); serviceName != "" {
		return serviceName
	}
	return appName
}