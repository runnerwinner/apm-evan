package api

import (
	"dogapm"
	"net/http"
	"ordersvc/grpcclient"
	"ordersvc/metric"
	"protos"
	"strconv"

	"github.com/google/uuid"
)

type order struct {
	
}

var Order = &order{}

func (o *order) Add(w http.ResponseWriter, r *http.Request) {

    ctx, span := dogapm.Tracer.Start(r.Context(), "orderStart.Add")
    defer span.End()

	// 获取参数
    values := r.URL.Query()

    uidStr := values.Get("uid")
    skuIDStr := values.Get("sku_id")
    numStr := values.Get("num")

    uid, err := strconv.Atoi(uidStr)
    if err != nil || uid <= 0 {
        http.Error(w, "invalid uid", http.StatusBadRequest)
        return
    }

    skuID, err := strconv.Atoi(skuIDStr)
    if err != nil || skuID <= 0 {
        http.Error(w, "invalid sku_id", http.StatusBadRequest)
        return
    }

    num, err := strconv.Atoi(numStr)
    if err != nil || num <= 0 {
        http.Error(w, "invalid num", http.StatusBadRequest)
        return
    }


	//检查用户信息
	_, err = grpcclient.UserClient.GetUserInfo(ctx, &protos.UserMsg{
		Id: int64(uid),
	})
	if err != nil {
		dogapm.Logger.Error(ctx, "get_userInfo", map[string]interface{}{
			"uid": uid,
		}, err)
		dogapm.HttpStatus.Error(w, err.Error(), nil)
		return
	}

	// 对库存进行扣减
	skuMsg, err := grpcclient.SkuClient.DecreaseStock(ctx, &protos.SkuMsg{
		Id:  int64(skuID),
		Num: int32(num),
	})
	if err != nil {
		dogapm.Logger.Error(ctx, "createOrder", map[string]interface{}{
			"uid":    uid,
			"sku_id": skuID,
			"num":    num,
		}, err)
		dogapm.HttpStatus.Error(w, err.Error(), nil)
		return
	}

	// 生成订单
	_, err = dogapm.Infra.Db.ExecContext(ctx, "INSERT INTO t_order (order_id, sku_id, num, price, uid) VALUES (?, ?, ?, ?, ?)", uuid.New().String(), skuID, num, skuMsg.Price, uid)
	if err != nil {
		dogapm.Logger.Error(ctx, "create_order_failed", map[string]interface{}{
			"uid":    uid,
			"sku_id": skuID,
		}, err)
		dogapm.HttpStatus.Error(w, err.Error(), nil)
		return
	}

	metric.OrderSuccessCounter.WithLabelValues(strconv.Itoa(skuID)).Inc()
	// 返回结果
	dogapm.HttpStatus.OK(w)

}