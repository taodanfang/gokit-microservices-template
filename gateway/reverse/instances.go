package reverse

import (
	"cctable/common/tools"
	"cctable/config"
	"github.com/hashicorp/consul/api"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Service_instance_reverse_manager struct {
	api_client *api.Client
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_service_instance_reverse_manager() *Service_instance_reverse_manager {

	cfg := config.Get_config()
	consul_config := api.DefaultConfig()
	consul_config.Address = cfg.Consul.Host + ":" + cfg.Consul.Port

	api_client, err := api.NewClient(consul_config)
	if err != nil {
		tools.Log(err)
		return nil
	}

	return &Service_instance_reverse_manager{
		api_client: api_client,
	}
}
