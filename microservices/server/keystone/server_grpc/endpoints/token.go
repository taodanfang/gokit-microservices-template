package endpoints

import (
	"cctable/api"
	"cctable/common/libs/micro/server/grpc_server/app"
	"cctable/common/results"
	"cctable/microservices/server/keystone/router/grpc_controller"
	"context"
	"github.com/go-kit/kit/endpoint"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Token_grpc_manager struct {
	transport_server app.I_grpc_application
	endpoints        map[string]endpoint.Endpoint

	ctl_token *grpc_controller.Token_grpc_controller
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_token_grpc_manager(tps app.I_grpc_application) {
	m := &Token_grpc_manager{transport_server: tps}

	m.ctl_token = grpc_controller.New_token_grpc_controller(tps)

	m.endpoints = make(map[string]endpoint.Endpoint, 0)

	m.Register_grpc_endpoint(api.EDP_grpc_keystone_token_check_access_token, m.check_access_token())
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (m *Token_grpc_manager) Register_grpc_endpoint(route_name string, edp endpoint.Endpoint) {
	tps := m.transport_server
	m.endpoints[route_name] = edp
	tps.Make_route(route_name, m.endpoints[route_name])
}

// --------------------------------------------------------------------
// Endpoint 方法
// --------------------------------------------------------------------

func (m *Token_grpc_manager) check_access_token() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		data := results.API()

		rq := request.(*app.Grpc_payload_request)

		err = rq.Check_request(
			api.EDP_grpc_keystone_token_check_access_token,
			"token_value")

		if err != nil {
			rs := results.Error(data, err.Error())
			return app.Grpc_payload_response{Result: *rs}, nil
		}

		token_value := rq.Params["token_value"].(string)

		rs := m.ctl_token.Check_access_token(ctx, token_value)

		return app.Grpc_payload_response{Result: *rs}, nil
	}
}
