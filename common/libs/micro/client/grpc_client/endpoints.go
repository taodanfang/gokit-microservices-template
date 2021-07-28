package grpc_client

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kataras/iris/v12"
	"strings"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_grpc_endpoint_client_manager interface {
	Register_grpc_endpoint(service_name, api_name, method string)
	Call_grpc_endpoint(ctx context.Context, service_name, method string, params iris.Map) *Json_rpc_response
}

type Grpc_endpoint_client_manager struct {
	transport_client I_grpc_client
	endpoints        map[string]endpoint.Endpoint
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func Make_grpc_endpoint_client_manager(tps I_grpc_client) *Grpc_endpoint_client_manager {

	endpoints := make(map[string]endpoint.Endpoint, 0)

	m := &Grpc_endpoint_client_manager{
		transport_client: tps,
		endpoints:        endpoints,
	}

	return m
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (c *Grpc_endpoint_client_manager) Register_grpc_endpoint(service_name, api_name, method string) {
	tps := c.transport_client
	edp := tps.Make_endpoint(service_name, method)
	c.endpoints[service_name+":"+api_name] = edp
	//log.Printf("edp: %#v\n", edp)
}

func (c *Grpc_endpoint_client_manager) Call_grpc_endpoint(ctx context.Context, service_name, api_name string, params iris.Map) *Json_rpc_response {

	api_name = strings.ToLower(api_name)

	//tools.Log(service_name, api_name)

	request := Json_rpc_request{
		Method: api_name,
		Params: params,
	}

	edp := c.endpoints[service_name+":"+api_name]

	//tools.Log(edp)

	rs, _ := edp(ctx, &request)

	return rs.(*Json_rpc_response)
}
