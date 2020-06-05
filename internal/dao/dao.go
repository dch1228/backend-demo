package dao

import (
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	"github.com/duchenhao/backend-demo/internal/conf"
	"github.com/duchenhao/backend-demo/internal/log"
	"github.com/duchenhao/backend-demo/internal/model"
)

var (
	db    *xorm.Engine
	cache *redis.Client
)

func Init() {
	initDB()

	initCache()
}

func initDB() {
	logger := log.Named("xorm")
	logger.Info("connecting to db", zap.String("dsn", conf.GetDBDsn()))

	var err error
	db, err = xorm.NewEngine("mysql", conf.GetDBDsn())
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxOpenConns(conf.DB.MaxOpenConn)
	db.DB().SetMaxIdleConns(conf.DB.MaxIdleConn)

	db.SetLogger(newLogger(logger))
	db.ShowSQL(true)

	db.SetMapper(names.GonicMapper{})

	if err := db.Ping(); err != nil {
		panic(err)
	}

}

func initCache() {
	cache = redis.NewClient(&redis.Options{
		Addr:     conf.Cache.Addr,
		Password: conf.Cache.Password,
		DB:       conf.Cache.DB,
	})

	_, err := cache.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func Migrate() {
	err := db.Sync2(
		model.User{},
	)
	if err != nil {
		panic(err)
	}
}

func Close() {
	db.Close()
	cache.Close()
}
