package naming

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func RegisterServiceInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) error {
	success, err := client.RegisterInstance(param)
	if err != nil {
		return err
	}
	fmt.Printf("RegisterServiceInstance, param:+%v,result:%+v \n", param, success)
	return nil
}

func Subscribe(client naming_client.INamingClient, param *vo.SubscribeParam) error {
	err := client.Subscribe(param)
	if err != nil {
		return err
	}
	return nil
}
