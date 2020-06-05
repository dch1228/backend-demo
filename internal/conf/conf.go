package conf

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type config struct {
	Core  *core  `mapstructure:"core"`
	DB    *db    `mapstructure:"db"`
	Cache *cache `mapstructure:"cache"`
}

type core struct {
	Name string `mapstructure:"name"`

	Addr  string `mapstructure:"addr"`
	Debug bool   `mapstructure:"debug"`

	MetricsUsername string `mapstructure:"metrics_username"`
	MetricsPassword string `mapstructure:"metrics_password"`

	Secret string `mapstructure:"secret"`
}

type db struct {
	Addr        string `mapstructure:"addr"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Database    string `mapstructure:"database"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
	MaxIdleConn int    `mapstructure:"max_Idle_conn"`
}

type cache struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

var (
	Core  = &core{}
	DB    = &db{}
	Cache = &cache{}
)

func Init(path string) {
	viper.SetConfigType("yaml")

	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := viper.ReadConfig(bytes.NewBuffer(b)); err != nil {
		panic(err)
	}

	cfg := &config{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}

	Core = cfg.Core
	DB = cfg.DB
	Cache = cfg.Cache
}

func GetHttpEnv() string {
	if Core.Debug {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}

func GetDBDsn() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", DB.User, DB.Password, DB.Addr, DB.Database)
	return dsn
}
