package configuration

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/zouyx/agollo/v4"
	"github.com/zouyx/agollo/v4/env/config"
	"io/ioutil"
	"strings"
)

func LoaderConfig(configPath, configName, configType string, configMapping interface{}) *viper.Viper {
	//data := &Application{}
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(configMapping); err != nil {
		panic(err)
	}
	return viper.GetViper()
}

func LoadConfigFromApollo(apolloConfigFile string) (*viper.Viper, error) {
	var c config.AppConfig
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		bt, err := ioutil.ReadFile(apolloConfigFile)
		if err != nil {
			panic(errors.WithMessage(err, "apollo文件读取异常"))
		}
		err = json.Unmarshal(bt, &c)
		return &c, err
	})
	if err != nil {
		panic(errors.WithMessage(err, "apollo启动失败"))
	}
	nss := strings.Split(c.NamespaceName, ",")
	v := viper.New()
	v.SetConfigType("prop")
	content := new(strings.Builder)
	for _, ns := range nss {
		nsConf := client.GetConfig(ns)
		if nsConf != nil {
			content.WriteString(nsConf.GetContent())
		}
	}
	v.ReadConfig(strings.NewReader(content.String()))
	return v, err
}
