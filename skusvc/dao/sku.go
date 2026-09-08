package dao

import (
	"context"
	"database/sql"
	"dogapm"
	"errors"
)

type skuDao struct{}

var SkuDao = &skuDao{}

func (s *skuDao) GetSkuByID(ctx context.Context, skuID int64) (map[string]interface{}, error) {
	info := dogapm.DBUtil.QueryFirst(
		dogapm.Infra.Db.QueryContext(
			ctx, "select * from t_sku where id = ?;", skuID,
		),
	)
	if info == nil {
		return nil, errors.New("sku not found")
	}
	return info, nil
}


func (s *skuDao) DecreaseStock(ctx context.Context, skuID int64, num int) (sql.Result,error) {
	result, err := dogapm.Infra.Db.ExecContext(ctx, "UPDATE t_sku SET num = num - ? WHERE id = ?", num, skuID)
	if err != nil {
		dogapm.Logger.Error(ctx, "decrease_stock", map[string]any{
			"sku_id": skuID,
			"num":    num,
		}, err)
		return nil, err
	}
	return result, nil
}