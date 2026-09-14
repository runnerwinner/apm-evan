package dao

import (
	"dogapm"
)

type deployInfo struct {}

var DeployInfo = &deployInfo{}

func (d *deployInfo) GetInfoByApp(app string) map[string]any {
	return dogapm.DBUtil.QueryFirst(dogapm.Infra.Db.Query(
		"select * from t_deploy_info where app = ?;",app))
}