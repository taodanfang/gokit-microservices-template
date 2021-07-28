package endpoints

import (
	"cctable/common/libs/micro/server/http_server/app"
	"cctable/common/libs/micro/server/http_server/transports"
	"cctable/common/results"
	"cctable/microservices/server/keystone/business/model"
	"cctable/microservices/server/keystone/business/service"
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
	New_user_http_manager(tps)
	New_oauth_http_manager(tps)
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func Check_oauth_authorization(ctx context.Context, app app.I_http_application) (client_uuid string, err error) {

	if err, ok := ctx.Value(transports.Authorization_error_key).(error); ok {
		return "", err
	}

	auth := ctx.Value(transports.Authorization_value_key).(iris.Map)
	the_auth_type := auth["auth_type"].(string)

	if the_auth_type == "oauth" {

		client_uuid = auth["client_uuid"].(string)
		client_secret := auth["client_secret"].(string)

		handler_client := service.New_client_service(app.Get_db())

		rs := handler_client.Check_client_by_uuid_and_secret(ctx, client_uuid, client_secret)
		if rs.Is_failure() {
			return "", results.Err_invalid_client_request
		}

		//tools.Log("client_uuid: ", client_uuid)

		return client_uuid, nil
	}

	return "", results.Err_invalid_authorization_header
}

func Check_token_authorization(ctx context.Context, app app.I_http_application) (token_value string, err error) {

	if err, ok := ctx.Value(transports.Authorization_error_key).(error); ok {
		return "", err
	}

	auth := ctx.Value(transports.Authorization_value_key).(iris.Map)
	the_auth_type := auth["auth_type"].(string)

	if the_auth_type == "token" {
		token_value = auth["token_value"].(string)

		handler_token := service.New_token_store(app.Get_db())
		rs := handler_token.Read_access_token(ctx, token_value)
		if rs.Is_failure() {
			return "", results.Err_token_does_not_exist
		}

		the_token := rs.Data["token"].(model.Token)

		if the_token.Is_expired() {
			return "", results.Err_token_has_expired
		}

		//tools.Log("token: ", token_value)

		return token_value, nil
	}

	return "", results.Err_invalid_authorization_header
}
