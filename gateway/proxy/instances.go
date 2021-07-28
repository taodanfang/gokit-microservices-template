package proxy

import (
	"cctable/common/tools"
	"cctable/config"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Service_instance_manager struct {
	client       consul.Client
	instance_map map[string]*consul.Instancer
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_service_instance_manager() *Service_instance_manager {

	cfg := config.Get_config()
	consul_config := api.DefaultConfig()
	consul_config.Address = cfg.Consul.Host + ":" + cfg.Consul.Port

	consul_client, err := api.NewClient(consul_config)
	if err != nil {
		tools.Log(err)
		return nil
	}
	client := consul.NewClient(consul_client)

	return &Service_instance_manager{
		client:       client,
		instance_map: make(map[string]*consul.Instancer, 0),
	}
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *Service_instance_manager) Register_instance(service_name string) *consul.Instancer {
	instance, ok := s.instance_map[service_name]
	if ok {
		return instance
	}

	instance = consul.NewInstancer(
		s.client, tools.Get_gokit_logger(), service_name, []string{"http"}, true)

	tools.Log("register_instance: ", instance)

	s.instance_map[service_name] = instance

	return instance
}

func (s *Service_instance_manager) Get_instance(service_name string) *consul.Instancer {
	instance, ok := s.instance_map[service_name]
	if ok {
		return instance
	}

	return nil
}
