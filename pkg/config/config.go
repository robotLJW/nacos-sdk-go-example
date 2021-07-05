package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once

var ConfigMessage *Config

func ReadConfig(configName string, configPath string, configType string) {
	once.Do(func() {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName(configName)
		viper.AddConfigPath(path + configPath)
		viper.SetConfigType(configType)
		err = viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		serverMessage := ServerMessage{
			IpAddr: viper.GetString("serverconfig.ipaddr"),
			Port:   viper.GetUint64("serverconfig.port"),
		}

		clientMessage := ClientMessage{
			NamespaceId: viper.GetString("clientconfig.namespaceId"),
			ServiceName: viper.GetString("clientconfig.serviceName"),
			LogDir:      viper.GetString("clientconfig.logDir"),
			CacheDir:    viper.GetString("clientconfig.cacheDir"),
			RotateTime:  viper.GetString("clientconfig.rotateTime"),
			MaxAge:      viper.GetInt64("clientconfig.maxAge"),
			LogLevel:    viper.GetString("clientconfig.logLevel"),
		}

		basicMessage := BasicMessage{
			InstanceIp:          viper.GetString("basicconfig.instanceIp"),
			InstancePort:        viper.GetUint64("basicconfig.instancePort"),
			InstanceCount:       viper.GetInt("basicconfig.instanceCount"),
			InstanceClusterName: viper.GetString("basicconfig.instanceClusterName"),
			SubscribeScope:      viper.GetInt("basicconfig.subscribeScope"),
		}
		ConfigMessage = &Config{
			Server: serverMessage,
			Client: clientMessage,
			Basic:  basicMessage,
		}
	})
}
