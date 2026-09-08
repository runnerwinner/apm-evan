package dogapm

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

type infra struct {
    Db *sql.DB
	Rdb *redis.Client
}

var Infra = &infra{}


type InfraOption func(i *infra)

func InfraDbOption(connectUrl string) InfraOption {
	return func(i *infra) {
		db, err := sql.Open("mysql", connectUrl)
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		i.Db = db
	}
}

func InfraRdbOption(connectUrl string) InfraOption {
	return func(i *infra) {
		rdb := redis.NewClient(&redis.Options{
			Addr: connectUrl,
			DB: 0,
		})
		if err := rdb.Ping(context.TODO()).Err(); err != nil {
			panic(err)
		}
		i.Rdb = rdb
	}
}

func (i *infra) Init(options ...InfraOption) {
	for _, option := range options {
		option(i)
	}
}