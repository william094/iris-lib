package iris_lib

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

var (
	Env    = os.Getenv("ENVIRON")
	GOCONF = os.Getenv("GOCONF")
)

func LoaderConfig(serverName string, rawValue *interface{}) {
	if Env == "" || Env == "develop" {
		Env = "dev"
	}
	if GOCONF == "" || GOCONF == "application" {
		GOCONF = "./etc/"
	}
	viper.AddConfigPath(GOCONF)
	viper.SetConfigName(serverName + "-" + Env)
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		SystemLogger.Error("配置文件加载失败", zap.Error(err))
		panic(err)
	}
	if err := viper.Unmarshal(rawValue); err != nil {
		SystemLogger.Error("配置文件加载失败", zap.Error(err))
		panic(err)
	}
}
