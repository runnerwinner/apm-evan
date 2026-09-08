package dao

import (
	"context"
	"dogapm"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type userDao struct{}

var UserDao = &userDao{}

func (s *userDao) GetUserByID(ctx context.Context, userID int64) (map[string]any, error) {

    // Try to get user info from Redis cache first
	usercache := dogapm.Infra.Rdb.Get(ctx, fmt.Sprintf("%s:%s:%d","usersvc","uid",userID))
	if usercache.Err() == nil {	
		userinfo := make(map[string]any)
		cacheStr, _ := usercache.Result()
		err := json.Unmarshal([]byte(cacheStr), &userinfo)
		if err == nil {
			return userinfo, nil
		}
	}

	// If not found in cache, query from database
	info := dogapm.DBUtil.QueryFirst(
		dogapm.Infra.Db.QueryContext(
			ctx, "select * from t_user where id = ?;", userID,
		),
	)
	if info == nil {
		return nil, errors.New("user not found")
	}

	cacheUserStr, err := json.Marshal(info)
	if err == nil {
		dogapm.Infra.Rdb.Set(ctx, fmt.Sprintf("%s:%s:%d","usersvc","uid",userID), cacheUserStr, 10*time.Minute)
	}
	return info, nil
}

