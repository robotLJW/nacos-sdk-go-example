package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"io/ioutil"
	"nacos-sdk-go-example/pkg/config"
	"net/http"

	"nacos-sdk-go-example/pkg/naming"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type User struct {
	Id   int
	Name string
	sex  int
}

func Hello(response http.ResponseWriter, request *http.Request) {

	user := &User{}
	user.Id = 1
	user.Name = "du"
	user.sex = 1
	data, _ := json.Marshal(user)
	fmt.Fprintf(response, string(data))
}

func GetJson(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	data, _ := ioutil.ReadAll(request.Body)
	user := &User{}
	_ = json.Unmarshal(data, user)

	user.Id = 2
	userJson, _ := json.Marshal(user)
	fmt.Fprintf(response, string(userJson))
}

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

	http.HandleFunc("/", Hello)
	http.HandleFunc("/getJson", GetJson)
	http.ListenAndServe(":8080", nil)
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
