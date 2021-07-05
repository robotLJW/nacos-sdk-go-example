package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"nacos-sdk-go-example/pkg/config"
	"nacos-sdk-go-example/pkg/naming"
	"nacos-sdk-go-example/pkg/uuid"
)

var wg sync.WaitGroup

func main() {
	config.ReadConfig("config", "/conf", "yaml")
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(config.ConfigMessage.Server.IpAddr, config.ConfigMessage.Server.Port),
	}
	cc := constant.ClientConfig{
		NamespaceId:         config.ConfigMessage.Client.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              config.ConfigMessage.Client.LogDir,
		CacheDir:            config.ConfigMessage.Client.CacheDir,
		RotateTime:          config.ConfigMessage.Client.RotateTime,
		MaxAge:              config.ConfigMessage.Client.MaxAge,
		LogLevel:            config.ConfigMessage.Client.LogLevel,
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	for err != nil {
		time.Sleep(10 * time.Second)
		client, err = clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig:  &cc,
				ServerConfigs: sc,
			},
		)
	}
	serviceName := uuid.GenerateServiceName()
	if len(config.ConfigMessage.Client.ServiceName) != 0 {
		serviceName = config.ConfigMessage.Client.ServiceName
	}
	instanceCount := config.ConfigMessage.Basic.InstanceCount

	for i := 1; i <= instanceCount; i++ {
		wg.Add(1)
		registerInstanceParam := vo.RegisterInstanceParam{
			Ip:          config.ConfigMessage.Basic.InstanceIp,
			Port:        config.ConfigMessage.Basic.InstancePort + uint64(i),
			ServiceName: serviceName,
			Weight:      10,
			ClusterName: config.ConfigMessage.Basic.InstanceClusterName,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
		}
		go registerInstance(client, registerInstanceParam)
	}
	wg.Wait()

	serviceNameOne := uuid.GenerateServiceName()

	registerInstanceParam := vo.RegisterInstanceParam{
		Ip:          config.ConfigMessage.Basic.InstanceIp,
		Port:        config.ConfigMessage.Basic.InstancePort,
		ServiceName: serviceNameOne,
		Weight:      10,
		ClusterName: config.ConfigMessage.Basic.InstanceClusterName,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}
	err = naming.RegisterServiceInstance(client, registerInstanceParam)
	for err != nil {
		fmt.Println(err)
		err = naming.RegisterServiceInstance(client, registerInstanceParam)
		if err == nil {
			break
		}
	}
	for i := 1; i <= 2; {
		service := randomServiceName("service", 500)
		err := naming.Subscribe(client, &vo.SubscribeParam{
			ServiceName: service,
			Clusters:    []string{config.ConfigMessage.Basic.InstanceClusterName},
			SubscribeCallback: func(instances []model.Instance, err error) {
				fmt.Printf("callback222 return instance:%+v \n", instances)
			},
		})
		for err != nil {
			fmt.Println(err)
			service = randomServiceName("service", 500)
			err = naming.Subscribe(client, &vo.SubscribeParam{
				ServiceName: service,
				Clusters:    []string{config.ConfigMessage.Basic.InstanceClusterName},
				SubscribeCallback: func(instances []model.Instance, err error) {
					fmt.Printf("callback222 return instance:%+v \n", instances)
				},
			})
		}
		i++
	}
	time.Sleep(360000 * time.Second)
}

func registerInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) {
	err := naming.RegisterServiceInstance(client, param)
	for {
		if err != nil {
			fmt.Println(err)
			err = naming.RegisterServiceInstance(client, param)
		} else {
			break
		}
	}
	wg.Done()
}

func randomServiceName(baseName string, scope int) string {
	randomNum := rand.Intn(scope)
	return baseName + strconv.Itoa(randomNum)
}
