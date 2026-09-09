package dogapm

import (
	"context"
	"testing"
)
func TestRedis(t *testing.T) {
	Infra.Init(
		InfraRdbOption("127.0.0.1:6380"),
		InfraEnableApm("127.0.0.1:54317"),
	)

	_, _ = Infra.Rdb.Get(context.TODO(), "test").Result()
	EndPoint.Close()
}