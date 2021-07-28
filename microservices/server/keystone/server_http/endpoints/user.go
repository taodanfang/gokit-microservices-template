package endpoints

import (
	"cctable/api"
	"cctable/common/libs/micro/server/http_server/app"
	"cctable/common/results"
	"cctable/microservices/server/keystone/router/http_controller"
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kataras/iris/v12"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type User_http_manager struct {
	transport_server app.I_http_application
	endpoints        map[string]endpoint.Endpoint

	ctl_user *http_controller.User_http_controller
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_user_http_manager(tps app.I_http_application) {
	m := &User_http_manager{transport_server: tps}

	m.ctl_user = http_controller.New_user_http_controller(tps)

	m.endpoints = make(map[string]endpoint.Endpoint, 0)

	m.Register_http_endpoint(api.EDP_http_keystone__user__register_user, m.register_user())
	m.Register_http_endpoint(api.EDP_http_keystone__user__login, m.login())
	m.Register_http_endpoint(api.EDP_http_keystone__user__get_all_users, m.get_all_users())
	m.Register_http_endpoint(api.EDP_http_keystone__user__logout, m.logout())
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (m *User_http_manager) Register_http_endpoint(route_name string, edp endpoint.Endpoint) {
	tps := m.transport_server
	m.endpoints[route_name] = edp
	tps.Make_route(route_name, m.endpoints[route_name])
}

// --------------------------------------------------------------------
// Endpoint 方法
// --------------------------------------------------------------------

func (m *User_http_manager) register_user() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		data := results.API()

		client_uuid, err := Check_oauth_authorization(ctx, m.transport_server)
		if err != nil {
			rs := results.Error(data, err.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		user_name := ""
		telephone := ""
		password := ""

		rq := request.(app.Http_payload_request)
		if value, ok := rq.Params["user_name"]; ok {
			user_name = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		if value, ok := rq.Params["telephone"]; ok {
			telephone = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		if value, ok := rq.Params["password"]; ok {
			password = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		rs := m.ctl_user.Register_user(ctx, client_uuid, user_name, telephone, password)
		if rs.Is_failure() {
			rs = results.Error(data, rs.Msg)
			return app.Http_payload_response{Result: *rs}, nil
		}

		data["result"] = iris.Map{
			"user": rs.Data["user"],
		}
		return app.Http_payload_response{Result: *results.Ok(data)}, nil
	}
}

func (m *User_http_manager) login() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		data := results.API()

		client_uuid, err := Check_oauth_authorization(ctx, m.transport_server)
		if err != nil {
			rs := results.Error(data, err.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		user_name := ""
		password := ""

		rq := request.(app.Http_payload_request)
		if value, ok := rq.Params["user_name"]; ok {
			user_name = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		if value, ok := rq.Params["password"]; ok {
			password = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		rs := m.ctl_user.Login(ctx, client_uuid, user_name, password)
		if rs.Is_failure() {
			rs = results.Error(data, rs.Msg)
			return app.Http_payload_response{Result: *rs}, nil
		}

		data["result"] = rs.Data
		return app.Http_payload_response{Result: *results.Ok(data)}, nil
	}
}

func (m *User_http_manager) get_all_users() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		data := results.API()

		_, err = Check_token_authorization(ctx, m.transport_server)
		if err != nil {
			rs := results.Error(data, err.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		rs := m.ctl_user.Get_all_users(ctx)
		if rs.Is_failure() {
			rs = results.Error(data, rs.Msg)
			return app.Http_payload_response{Result: *rs}, nil
		}

		data["result"] = rs.Data
		return app.Http_payload_response{Result: *results.Ok(data)}, nil
	}
}

func (m *User_http_manager) logout() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		data := results.API()

		token_value, err := Check_token_authorization(ctx, m.transport_server)

		if err != nil {
			rs := results.Error(data, err.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		user_uuid := ""

		rq := request.(app.Http_payload_request)
		if value, ok := rq.Params["user_uuid"]; ok {
			user_uuid = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		rs := m.ctl_user.Logout(ctx, token_value, user_uuid)
		if rs.Is_failure() {
			rs = results.Error(data, rs.Msg)
			return app.Http_payload_response{Result: *rs}, nil
		}

		data["result"] = rs.Data
		return app.Http_payload_response{Result: *results.Ok(data)}, nil
	}
}
