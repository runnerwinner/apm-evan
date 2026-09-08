package grpc

import (
	"context"
	"protos"
	"usrsvc/dao"

	"github.com/spf13/cast"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type UserServer struct {
	protos.UnimplementedUserServiceServer
}

func (s *UserServer) GetUserInfo(ctx context.Context, userIn *protos.UserMsg) (*protos.UserMsg, error) {
	//获取用户信息
	info, err := dao.UserDao.GetUserByID(ctx, userIn.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	if len(info) == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &protos.UserMsg{
		Name: cast.ToString(info["name"]),
		Id:   cast.ToInt64(info["id"]),
	}, nil
}
	