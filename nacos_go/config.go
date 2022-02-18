package nacos_go

import (
	"fmt"
	"log"

	// "github.com/laixyz/utils"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type ServiceParam vo.RegisterInstanceParam

type ConfigParam vo.ConfigParam

func ClientRegisterInstance(nacos []string, svc ServiceParam) {

	namingClient, _ := getNamingClient(nacos)
	// create namingClient Register Param
	param := vo.RegisterInstanceParam{
		Ip:          util.LocalIP(),
		Port:        8080,
		ServiceName: svc.ServiceName,
		GroupName:   svc.GroupName,
		ClusterName: svc.ClusterName,
		Weight:      svc.Weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"version": "v0.9.0", "env": "prod"},
	}

	success, _ := namingClient.RegisterInstance(param)
	fmt.Printf("ClientRegisterInstance, param:%+v,result:%+v \n\n", param, success)
}

func ClientSubscribeService(nacos []string, svc ServiceParam) {

	namingClient, _ := getNamingClient(nacos)
	// create namingClient Subscribe Param
	param := vo.SubscribeParam{
		ServiceName: svc.ServiceName,
		Clusters:    []string{svc.ClusterName},
		GroupName:   svc.GroupName,
		SubscribeCallback: func(services []model.Instance, err error) {
			log.Printf("\n\n callback return services:%s \n\n", util.ToJsonString(services))
		},
	}

	if err := namingClient.Subscribe(&param); err != nil {
		log.Printf("namingClient Subscribe failed [%v]", err)
	}

}

func ClientUnsubscribeService(nacos []string, svc ServiceParam) {
	namingClient, _ := getNamingClient(nacos)
	// create namingClient Unsubscribe Param
	// create namingClient Subscribe Param
	param := vo.SubscribeParam{
		ServiceName: svc.ServiceName,
		Clusters:    []string{svc.ClusterName},
		GroupName:   svc.GroupName,
		SubscribeCallback: func(services []model.Instance, err error) {
			log.Printf("\n\n callback return services:%s \n\n", util.ToJsonString(services))
		},
	}
	if err := namingClient.Unsubscribe(&param); err != nil {
		log.Printf("namingClient Unsubscribe failed [%v]", err)
	}
}

func ClientPublishConfig(nacos []string, cfg ConfigParam) {
	configClient, _ := getConfigClient(nacos)

	// create configClient Param
	param := vo.ConfigParam{
		DataId:  cfg.DataId,
		Group:   cfg.Group,
		Content: cfg.Content,
	}
	success, _ := configClient.PublishConfig(param)

	fmt.Printf("ClientPublishConfig, param:%+v,result:%+v \n\n", param, success)
}

func getNacosConfig(nacos []string) (clientConfig constant.ClientConfig, serverConfigs []constant.ServerConfig) {
	//create ServerConfig
	for _, addr := range nacos {
		serverConfigs = append(serverConfigs, *constant.NewServerConfig(addr, 8848))
	}

	//create ClientConfig
	clientConfig = *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("info"),
	)

	return clientConfig, serverConfigs
}

func getNamingClient(nacos []string) (namingClient naming_client.INamingClient, err error) {
	clientConfig, serverConfigs := getNacosConfig(nacos)
	// create naming client
	namingClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		panic(err)
	}

	return namingClient, err
}

func getConfigClient(nacos []string) (configClient config_client.IConfigClient, err error) {
	clientConfig, serverConfigs := getNacosConfig(nacos)
	// create config client
	configClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		panic(err)
	}

	return configClient, err
}
