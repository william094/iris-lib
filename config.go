package iris_lib

import (
	"github.com/spf13/viper"
	"os"
)

var (
	Env    = os.Getenv("ENVIRON")
	GOCONF = os.Getenv("GOCONF")
)

func LoaderConfig(serverName string) (*Application, *viper.Viper) {
	data := &Application{}
	if Env == "" || Env == "develop" {
		Env = "dev"
	}
	if GOCONF == "" {
		GOCONF = "./etc/"
	}
	viper.AddConfigPath(GOCONF)
	viper.SetConfigName(serverName + "-" + Env)
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(data); err != nil {
		panic(err)
	}
	return data, viper.GetViper()
}
