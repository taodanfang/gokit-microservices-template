package endpoints

import (
	"cctable/common/libs/micro/server/http_server/app"
	"cctable/common/libs/micro/server/http_server/transports"
	"cctable/common/results"
	"cctable/common/tools"
	"cctable/microservices/client/keystone/client_grpc/controller"
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kataras/iris/v12"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_http_endpoint_manager interface {
	Register_http_endpoint(route_name string, edp endpoint.Endpoint)
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func Init_http_endpoint_managers(tps *transports.Transport_for_http_application) {
	New_device_http_manager(tps)
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func Check_token_authorization(ctx context.Context, app app.I_http_application) error {
	rpc, err := controller.New_keystone_token_grpc_client_controller()
	defer rpc.Close_client()

	if err != nil {
		return err
	}

	if err, ok := ctx.Value(transports.Authorization_error_key).(error); ok {
		return err
	}

	auth := ctx.Value(transports.Authorization_value_key).(iris.Map)
	the_auth_type := auth["auth_type"].(string)

	if the_auth_type == "token" {
		token_value := auth["token_value"].(string)
		//tools.Log("rpc.check_token: ", token_value)

		rs := rpc.Check_access_token(ctx, token_value)

		tools.Log(rs)

		if rs.Success == false {
			return results.Err_token_does_not_exist
		}

		return nil
	}

	return results.Err_invalid_authorization_header
}
