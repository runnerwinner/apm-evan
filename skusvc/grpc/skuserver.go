package grpc

import (
	"context"
	"protos"
	"skusvc/dao"

	"github.com/spf13/cast"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type SkuServer struct {
	protos.UnimplementedSkuServiceServer
}

func (s *SkuServer) DecreaseStock(ctx context.Context, skuIn *protos.SkuMsg) (*protos.SkuMsg, error) {
	//获取商品信息
	info, err := dao.SkuDao.GetSkuByID(ctx, skuIn.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	if len(info)==0 {
		return nil, status.Errorf(codes.NotFound, "sku not found")
	}
	//进行扣减库存
	decrRes, err := dao.SkuDao.DecreaseStock(ctx, skuIn.Id, int(skuIn.Num))
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	if rows, _ := decrRes.RowsAffected(); rows == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "decrease stock failed")
	}

	return &protos.SkuMsg{
		Name:  cast.ToString(info["name"]),
		Id:    cast.ToInt64(info["id"]),
		Price: cast.ToInt32(info["price"]),
		Num:   cast.ToInt32(info["num"]) - skuIn.Num,
	}, nil
}