package discover_http

import (
	"context"
	"errors"
)

// --------------------------------------------------------------------
// 错误处理
// --------------------------------------------------------------------

var (
	Err_service_instance_does_not_exist = errors.New("错误：Service 实例不存在")
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_http_discovery_service interface {
	Health_check() bool
	Discovery_service(ctx context.Context, service_name string) ([]interface{}, error)
}

type Http_discovery_service struct {
	transport_server I_http_discovery
	discovery_client I_discovery_http_client
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_http_discovery_service(tps I_http_discovery) I_http_discovery_service {
	return &Http_discovery_service{
		transport_server: tps,
		discovery_client: tps.Get_discovery_client(),
	}
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *Http_discovery_service) Discovery_service(_ context.Context, service_name string) ([]interface{}, error) {

	instances := s.discovery_client.Discover_services(service_name)

	if instances == nil || len(instances) == 0 {
		return nil, Err_service_instance_does_not_exist
	}

	return instances, nil
}

func (s *Http_discovery_service) Health_check() bool {
	return true
}
