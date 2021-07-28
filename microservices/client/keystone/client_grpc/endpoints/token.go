package endpoints

import (
	"cctable/api"
	"cctable/common/libs/micro/client/grpc_client"
	"context"
	"github.com/kataras/iris/v12"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_token_service_grpc_client interface {
	Check_access_token(ctx context.Context, token_value string) *grpc_client.Json_rpc_response
}

type Token_service_grpc_client struct {
	manager      grpc_client.I_grpc_endpoint_client_manager
	service_name string
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_token_service_grpc_client(manager grpc_client.I_grpc_endpoint_client_manager) I_token_service_grpc_client {
	s := &Token_service_grpc_client{
		manager:      manager,
		service_name: api.SVC_grpc_keystone_pb_token_service,
	}

	manager.Register_grpc_endpoint(
		s.service_name,
		api.EDP_grpc_keystone_token_check_access_token, "CheckAccessToken")
	return s
}

// --------------------------------------------------------------------
// RPC 可调用方法
// --------------------------------------------------------------------

func (c *Token_service_grpc_client) Check_access_token(ctx context.Context, token_value string) *grpc_client.Json_rpc_response {

	params := iris.Map{
		"token_value": token_value,
	}

	//tools.Log(params)

	rs := c.manager.Call_grpc_endpoint(
		ctx,
		c.service_name,
		api.EDP_grpc_keystone_token_check_access_token, params)

	//tools.Log(rs)

	return rs
}
