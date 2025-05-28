package registry

import (
	"layout/internal/conf"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	nacos "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
)

func NewNacosRegistry(c *conf.Registry) registry.Registrar {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(c.GetNacos().GetAddr(), c.GetNacos().GetPort()),
	}
	cc := constant.NewClientConfig(constant.WithNamespaceId(c.GetNacos().GetNamespaceId()))
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ServerConfigs: sc,
		ClientConfig:  cc,
	})
	if err != nil {
		panic(err)
	}
	return nacos.New(client)
}
