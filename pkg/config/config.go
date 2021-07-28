package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once

var ConfigMessage *Config

func ReadConfig(configName string, configPath string, configType string) {
	once.Do(func() {
		viper.SetConfigName(configName)
		viper.AddConfigPath(configPath)
		viper.SetConfigType(configType)
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		serverMessage := ServerMessage{
			IpAddr: viper.GetString("serverConfig.ipaddr"),
			Port:   viper.GetUint64("serverConfig.port"),
		}

		clientMessage := ClientMessage{
			NamespaceId: viper.GetString("clientConfig.namespaceId"),
			ServiceName: viper.GetString("clientConfig.serviceName"),
			Scope:       viper.GetInt("clientConfig.scope"),
			LogDir:      viper.GetString("clientConfig.logDir"),
			CacheDir:    viper.GetString("clientConfig.cacheDir"),
			RotateTime:  viper.GetString("clientConfig.rotateTime"),
			MaxAge:      viper.GetInt64("clientConfig.maxAge"),
			LogLevel:    viper.GetString("clientConfig.logLevel"),
		}

		basicMessage := BasicMessage{
			InstanceIp:             viper.GetString("basicConfig.instanceIp"),
			InstancePort:           viper.GetUint64("basicConfig.instancePort"),
			SubscribeInstanceCount: viper.GetInt("basicConfig.subscribeInstanceCount"),
			InstanceClusterName:    viper.GetString("basicConfig.instanceClusterName"),
			SubscribeScope:         viper.GetInt("basicConfig.subscribeScope"),
			NameServerAddr:         viper.GetString("basicConfig.nameServerAddr"),
		}
		ConfigMessage = &Config{
			Server: serverMessage,
			Client: clientMessage,
			Basic:  basicMessage,
		}
	})
}
