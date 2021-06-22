package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var (
	once sync.Once

	serverIpAddr string
	serverPort   uint64

	clientNamespaceId string
	clientLogDir      string
	clientCacheDir    string
	clientRotateTime  string
	clientMaxAge      int64
	clientLogLevel    string
)

func main() {
	readConfig()
	sc := []constant.ServerConfig{
		{
			IpAddr: serverIpAddr,
			Port:   serverPort,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         clientNamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              clientLogDir,
		CacheDir:            clientCacheDir,
		RotateTime:          clientRotateTime,
		MaxAge:              clientMaxAge,
		LogLevel:            clientLogLevel,
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	serviceName, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	registerInstanceParam := vo.RegisterInstanceParam{
		Ip:          "10.10.10.11",
		Port:        8848,
		ServiceName: serviceName.String(),
		Weight:      10,
		ClusterName: "cluster-a",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}
	err = registerServiceInstance(client, registerInstanceParam)
	for {
		if err != nil {
			err = registerServiceInstance(client, registerInstanceParam)
		} else {
			break
		}
	}
	time.Sleep(360000 * time.Second)
}

func readConfig() {
	once.Do(func() {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("config")
		viper.AddConfigPath(path + "/discovery/conf")
		viper.SetConfigType("yaml")
		err = viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		serverIpAddr = viper.GetString("serverconfig.ipaddr")
		serverPort = viper.GetUint64("serverconfig.port")

		clientNamespaceId = viper.GetString("clientconfig.namespaceId")
		clientLogDir = viper.GetString("clientconfig.logDir")
		clientCacheDir = viper.GetString("clientconfig.cacheDir")
		clientRotateTime = viper.GetString("clientconfig.rotateTime")
		clientMaxAge = viper.GetInt64("clientconfig.maxAge")
		clientLogLevel = viper.GetString("clientconfig.logLevel")
	})
}

func registerServiceInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) error {
	success, err := client.RegisterInstance(param)
	if err != nil {
		return err
	}
	fmt.Printf("RegisterServiceInstance, param:+%v,result:%+v \n", param, success)
	return nil
}
