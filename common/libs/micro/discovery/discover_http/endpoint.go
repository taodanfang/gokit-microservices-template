package discover_http

import (
	"github.com/go-kit/kit/endpoint"
	"net/http"

	"context"
	"errors"
)

// --------------------------------------------------------------------
// 错误处理
// --------------------------------------------------------------------

var (
	Err_bad_request = errors.New("错误：无效请求")
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Http_discovery_controller struct {
	transport_server I_http_discovery
	discovery_client I_discovery_http_client

	Discovery_endpoint    endpoint.Endpoint
	Health_check_endpoint endpoint.Endpoint

	handler_discovery I_http_discovery_service
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_http_discovery_controller(tps I_http_discovery) {
	ctl := &Http_discovery_controller{
		transport_server: tps,
		discovery_client: tps.Get_discovery_client()}

	ctl.handler_discovery = New_http_discovery_service(tps)
	ctl.Discovery_endpoint = ctl.Make_discovery_endpoint()
	ctl.Health_check_endpoint = ctl.Make_health_check_endpoint()

	ctl.register_endpoint_router()
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (ctl *Http_discovery_controller) register_endpoint_router() {
	tps := ctl.transport_server

	tps.Make_route("discovery", ctl.Discovery_endpoint, decode_discovery_request)
	tps.Make_route("health", ctl.Health_check_endpoint, decode_health_check_request)
}

// --------------------------------------------------------------------
// Endpoint 方法(1)
// --------------------------------------------------------------------

// 服务发现请求结构体
type DiscoveryRequest struct {
	Service_name string
}

// 服务发现响应结构体
type DiscoveryResponse struct {
	Instances []interface{} `json:"instances"`
	Error     string        `json:"error"`
}

func decode_discovery_request(_ context.Context, r *http.Request) (interface{}, error) {
	service_name := r.URL.Query().Get("service_name")
	if service_name == "" {
		return nil, Err_bad_request
	}
	return DiscoveryRequest{
		Service_name: service_name,
	}, nil
}

// 创建服务发现的 Endpoint
func (ctl *Http_discovery_controller) Make_discovery_endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DiscoveryRequest)
		instances, err := ctl.handler_discovery.Discovery_service(ctx, req.Service_name)
		var errString = ""
		if err != nil {
			errString = err.Error()
		}
		return &DiscoveryResponse{
			Instances: instances,
			Error:     errString,
		}, nil
	}
}

// --------------------------------------------------------------------
// Endpoint 方法(2)
// --------------------------------------------------------------------

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

func decode_health_check_request(ctx context.Context, r *http.Request) (interface{}, error) {
	return HealthRequest{}, nil
}

// 创建健康检查的 Endpoint
func (ctl *Http_discovery_controller) Make_health_check_endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := ctl.handler_discovery.Health_check()
		return HealthResponse{
			Status: status,
		}, nil
	}
}
