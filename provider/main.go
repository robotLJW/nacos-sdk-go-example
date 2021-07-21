package main

import (
	"fmt"
	"math/rand"
	"nacos-sdk-go-example/pkg/name"
	"strconv"
	"strings"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"nacos-sdk-go-example/pkg/config"
	"nacos-sdk-go-example/pkg/naming"
)

func main() {
	config.ReadConfig("config", "/conf", "yaml")
	serviceNameAddr := config.ConfigMessage.Basic.NameServerAddr
	serviceName, err := name.ReadName(serviceNameAddr)
	if err != nil {
		panic(err)
	}
	ipAddr := strings.Split(config.ConfigMessage.Server.IpAddr, ",")
	sc := make([]constant.ServerConfig, 0)
	for i := 0; i < len(ipAddr); i++ {
		tmpServerConfig := *constant.NewServerConfig(ipAddr[i], config.ConfigMessage.Server.Port)
		sc = append(sc, tmpServerConfig)
		fmt.Println(ipAddr[i])
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
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	scope := config.ConfigMessage.Client.Scope

	instanceParam := vo.RegisterInstanceParam{
		Ip:          config.ConfigMessage.Basic.InstanceIp + strconv.Itoa(randomNumb(scope)),
		Port:        config.ConfigMessage.Basic.InstancePort + uint64(randomNumb(scope)),
		ServiceName: serviceName,
		Weight:      10,
		ClusterName: config.ConfigMessage.Basic.InstanceClusterName,
		GroupName:   "group-A",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}

	go registerInstance(client, instanceParam)

	time.Sleep(360000 * time.Second)

}

func randomNumb(scope int) int {
	return rand.Intn(scope)
}

func registerInstanceByScope(client naming_client.INamingClient, serviceName string, scope int) {
	registerInstanceParam := vo.RegisterInstanceParam{
		Ip:          config.ConfigMessage.Basic.InstanceIp,
		Port:        config.ConfigMessage.Basic.InstancePort,
		ServiceName: randomService(serviceName, scope),
		Weight:      10,
		ClusterName: config.ConfigMessage.Basic.InstanceClusterName,
		GroupName:   "group-A",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}
	fmt.Println(registerInstanceParam.ServiceName)
	err := naming.RegisterServiceInstance(client, registerInstanceParam)
	for {
		if err != nil {
			fmt.Println(err)
			registerInstanceParam.ServiceName = randomService(serviceName, scope)
			fmt.Println(registerInstanceParam.ServiceName)
			err = naming.RegisterServiceInstance(client, registerInstanceParam)
		} else {
			break
		}
	}
	svcName := registerInstanceParam.ServiceName
	instanceCount := config.ConfigMessage.Basic.InstanceCount
	for i := 1; i < instanceCount; i++ {
		instanceParam := vo.RegisterInstanceParam{
			Ip:          config.ConfigMessage.Basic.InstanceIp,
			Port:        config.ConfigMessage.Basic.InstancePort + uint64(i),
			ServiceName: svcName,
			Weight:      10,
			ClusterName: config.ConfigMessage.Basic.InstanceClusterName,
			GroupName:   "group-A",
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
		}
		go registerInstance(client, instanceParam)
	}

}

func randomService(serviceName string, scope int) string {
	return serviceName + strconv.Itoa(rand.Intn(scope))
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
}
