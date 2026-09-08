package dogapm

import (
	"context"
	"maps"

	"github.com/sirupsen/logrus"
)
type log struct{}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

var Logger = &log{}

func (l *log) buildFields(action string, kv map[string]interface{}) map[string]interface{} {
	fields := make(map[string]interface{}, len(kv)+1)
	maps.Copy(fields, kv)
	fields["action"] = action
	return fields
}

func (l *log) Info(ctx context.Context, action string, kv map[string]interface{}) {
	logrus.WithFields(l.buildFields(action, kv)).Info()
}

func (l *log) Warn(ctx context.Context, action string, kv map[string]interface{}) {
	logrus.WithFields(l.buildFields(action, kv)).Warn()
}

func (l *log) Debug(ctx context.Context, action string, kv map[string]interface{}) {
	logrus.WithFields(l.buildFields(action, kv)).Debug()
}

func (l *log) Error(ctx context.Context, action string, kv map[string]interface{}, err error) {
	fields := l.buildFields(action, kv)
	logrus.WithFields(fields).WithError(err).Error()
}