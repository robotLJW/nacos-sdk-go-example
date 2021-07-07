package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"io/ioutil"
	"nacos-sdk-go-example/pkg/config"
	"nacos-sdk-go-example/pkg/naming"
	"net/http"
)

var address string

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
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	serviceName := config.ConfigMessage.Client.ServiceName
	registerInstanceParam := vo.RegisterInstanceParam{
		Ip:          config.ConfigMessage.Basic.InstanceIp,
		Port:        config.ConfigMessage.Basic.InstancePort,
		ServiceName: serviceName,
		Weight:      10,
		ClusterName: config.ConfigMessage.Basic.InstanceClusterName,
		GroupName:   "group-A",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}
	go registerInstance(client, registerInstanceParam)

	instanceParam := vo.SelectOneHealthInstanceParam{
		Clusters:    []string{config.ConfigMessage.Basic.InstanceClusterName},
		ServiceName: "provider",
		GroupName:   "group-A",
	}
	instance, err := naming.GetOneHealthInstance(client, instanceParam)
	if err != nil {
		panic(err)
	}
	address = fmt.Sprintf("http://%s:%v/", instance.Ip, instance.Port)
	fmt.Println(address)

	err = naming.Subscribe(client, &vo.SubscribeParam{
		ServiceName: "provider",
		Clusters:    []string{config.ConfigMessage.Basic.InstanceClusterName},
		GroupName:   "group-A",
		SubscribeCallback: func(instances []model.Instance, err error) {
			if len(instances) != 0 {
				address = fmt.Sprintf("http://%s:%v/", instances[0].Ip, instances[0].Port)
			} else {
				address = ""
			}
			fmt.Println(fmt.Sprintf("callback service: %+v return instance:%+v \n", "provider", instances))
		},
	})
	http.HandleFunc("/", Hello)
	http.ListenAndServe(":8888", nil)
}

func getHelloFromProvider(address string) string {
	resp, _ := http.Get(address)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func Hello(response http.ResponseWriter, request *http.Request) {
	data := getHelloFromProvider(address)
	data = "from provider: " + data
	fmt.Fprintf(response, data)
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