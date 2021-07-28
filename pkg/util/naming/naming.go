package naming

import (
	"fmt"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"nacos-sdk-go-example/pkg/config"
)

func GetServerConfigs(ipAddresses string, port uint64) []constant.ServerConfig {
	ipAddr := strings.Split(ipAddresses, ",")
	sc := make([]constant.ServerConfig, 0)
	for i := 0; i < len(ipAddr); i++ {
		tmpServerConfig := *constant.NewServerConfig(ipAddr[i], port)
		sc = append(sc, tmpServerConfig)
	}
	return sc
}

func GetClientConfig() constant.ClientConfig {
	if config.ConfigMessage == nil {
		return constant.ClientConfig{}
	}
	return constant.ClientConfig{
		NamespaceId:         config.ConfigMessage.Client.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              config.ConfigMessage.Client.LogDir,
		CacheDir:            config.ConfigMessage.Client.CacheDir,
		RotateTime:          config.ConfigMessage.Client.RotateTime,
		MaxAge:              config.ConfigMessage.Client.MaxAge,
		LogLevel:            config.ConfigMessage.Client.LogLevel,
	}
}

func NewNamingClient(serverConfigs []constant.ServerConfig, clientConfig constant.ClientConfig) (naming_client.INamingClient, error) {
	return clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
}

func RegisterServiceInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) error {
	success, err := client.RegisterInstance(param)
	if err != nil {
		return err
	}
	fmt.Printf("RegisterServiceInstance, param:+%v,result:%+v \n", param, success)
	return nil
}

func GetOneHealthInstance(client naming_client.INamingClient, param vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	return client.SelectOneHealthyInstance(param)
}

func Subscribe(client naming_client.INamingClient, param *vo.SubscribeParam) error {
	err := client.Subscribe(param)
	if err != nil {
		return err
	}
	return nil
}
