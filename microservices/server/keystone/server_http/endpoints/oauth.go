package endpoints

import (
	"cctable/api"
	"cctable/common/libs/micro/server/http_server/app"
	"cctable/common/results"
	"cctable/microservices/server/keystone/business/helper"
	"cctable/microservices/server/keystone/router/http_controller"
	"context"
	"github.com/kataras/iris/v12"

	"github.com/go-kit/kit/endpoint"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 3.2) 实现服务 manager

type Oauth_http_manager struct {
	transport_server app.I_http_application
	endpoints        map[string]endpoint.Endpoint

	ctl_oauth *http_controller.Oauth_http_controller
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

// 3.3) 创建服务 manager

func New_oauth_http_manager(tps app.I_http_application) {
	m := &Oauth_http_manager{transport_server: tps}
	m.ctl_oauth = http_controller.New_oauth_http_controller(tps)

	m.endpoints = make(map[string]endpoint.Endpoint, 0)

	m.Register_http_endpoint(api.EDP_http_keystone__oauth__grant_token, m.grant_token())
	m.Register_http_endpoint(api.EDP_http_keystone__oauth__check_token, m.check_token())
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (m *Oauth_http_manager) Register_http_endpoint(route_name string, edp endpoint.Endpoint) {
	tps := m.transport_server
	m.endpoints[route_name] = edp
	tps.Make_route(route_name, m.endpoints[route_name])
}

// --------------------------------------------------------------------
// Endpoint 方法
// --------------------------------------------------------------------

// 3.5) 定义服务 endpoint

func (m *Oauth_http_manager) grant_token() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		data := results.API()

		client_uuid, err := Check_oauth_authorization(ctx, m.transport_server)
		if err != nil {
			rs := results.Error(data, err.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		grant_type := ""
		user_name := ""
		user_password := ""
		token_value := ""

		rq := request.(app.Http_payload_request)
		if value, ok := rq.Params["grant_type"]; ok {
			grant_type = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		if grant_type == helper.TOKEN_GRANT_type_username_and_password {

			if value, ok := rq.Params["user_name"]; ok {
				user_name = value.(string)
			} else {
				rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
				return app.Http_payload_response{Result: *rs}, nil
			}

			if value, ok := rq.Params["password"]; ok {
				user_password = value.(string)
			} else {
				rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
				return app.Http_payload_response{Result: *rs}, nil
			}

			rs := m.ctl_oauth.Grant_token_with_user_and_password(ctx, client_uuid, user_name, user_password)
			if rs.Is_failure() {
				rs = results.Error(data, rs.Msg)
				return app.Http_payload_response{Result: *rs}, nil
			}

			data["result"] = iris.Map{
				"access_token":  rs.Data["access_token"],
				"refresh_token": rs.Data["refresh_token"],
			}
			return app.Http_payload_response{Result: *results.Ok(data)}, nil

		}

		if grant_type == helper.TOKEN_GRANT_type_refresh_token {

			if value, ok := rq.Params["token_value"]; ok {
				token_value = value.(string)
			} else {
				rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
				return app.Http_payload_response{Result: *rs}, nil
			}

			rs := m.ctl_oauth.Grant_token_with_refresh_token(ctx, client_uuid, token_value)
			if rs.Is_failure() {
				rs = results.Error(data, rs.Msg)
				return app.Http_payload_response{Result: *rs}, nil
			}

			data["result"] = iris.Map{
				"access_token":  rs.Data["access_token"],
				"refresh_token": rs.Data["refresh_token"],
			}
			return app.Http_payload_response{Result: *results.Ok(data)}, nil
		}

		rs := results.Error(data, results.Err_does_not_support_the_grant_type.Error())
		return app.Http_payload_response{Result: *rs}, nil
	}
}

func (m *Oauth_http_manager) check_token() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		data := results.API()

		client_uuid, err := Check_oauth_authorization(ctx, m.transport_server)
		if err != nil {
			rs := results.Error(data, err.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		token_value := ""

		rq := request.(app.Http_payload_request)
		if value, ok := rq.Params["token_value"]; ok {
			token_value = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		rs := m.ctl_oauth.Check_token(ctx, client_uuid, token_value)
		if rs.Is_failure() {
			rs = results.Error(data, rs.Msg)
			return app.Http_payload_response{Result: *rs}, nil
		}

		data["result"] = iris.Map{
			"client": rs.Data["client"],
			"user":   rs.Data["user"],
		}
		return app.Http_payload_response{Result: *results.Ok(data)}, nil
	}
}
