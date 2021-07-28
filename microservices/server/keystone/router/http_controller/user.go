package http_controller

import (
	"cctable/common/libs/micro/server/http_server/app"
	"cctable/common/results"
	"cctable/common/tools"
	"cctable/microservices/server/keystone/business/helper"
	"cctable/microservices/server/keystone/business/model"
	"cctable/microservices/server/keystone/business/service"
	"context"
	"github.com/kataras/iris/v12"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type User_http_controller struct {
	server app.I_http_application

	handler_client      service.I_client_service
	handler_oauth       service.I_oauth_service
	handler_token_store service.I_token_store
	handler_user        service.I_user_service

	helper_token_enhance  helper.I_token_enhance_helper
	helper_token          helper.I_token_helper
	helper_token_granters helper.I_token_grant_helper
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_user_http_controller(server app.I_http_application) *User_http_controller {

	ctl := &User_http_controller{server: server}

	ctl.handler_client = service.New_client_service(server.Get_db())
	ctl.handler_oauth = service.New_oauth_service(server.Get_db())
	ctl.handler_token_store = service.New_token_store(server.Get_db())
	ctl.handler_user = service.New_user_service(server.Get_db())

	ctl.helper_token_enhance = helper.New_jwt_token_enhancer(server.Get_db(), "jwt_secret")
	ctl.helper_token = helper.New_token_helper(server.Get_db(), ctl.helper_token_enhance)
	token_granter_with_user_and_password := helper.New_token_granter_with_username_and_password(server.Get_db(), ctl.helper_token)
	ctl.helper_token_granters = helper.New_token_granters(map[string]helper.I_token_grant_helper{
		helper.TOKEN_GRANT_type_username_and_password: token_granter_with_user_and_password,
	})

	return ctl
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (ctl *User_http_controller) Register_user(ctx context.Context, client_uuid string, user_name, telephone, password string) *results.Result {

	data := results.API()

	rs := ctl.handler_user.Check_user_by_name_and_telephone(ctx, user_name, telephone)
	if rs.Is_success() {
		return results.Ok(data)
	}

	rs = ctl.handler_user.New_user_with_telephone(ctx, user_name, telephone)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_user := rs.Data["user"].(model.User)

	rs = ctl.handler_user.Update_user(ctx, the_user.User_uuid, "password", password)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_user = rs.Data["user"].(model.User)

	rs = ctl.handler_oauth.New_oauth(ctx, the_user.User_uuid, client_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	data["user"] = the_user

	return results.Ok(data)
}

func (ctl *User_http_controller) Login(ctx context.Context, client_uuid, user_name, password string) *results.Result {

	data := results.API()

	rs := ctl.handler_client.Get_client_by_uuid(ctx, client_uuid)
	if rs.Is_failure() {
		return results.Error(data, results.Err_client_does_not_exist.Error())
	}

	rs = ctl.handler_user.Check_user_by_name_and_password(ctx, user_name, password)
	if rs.Is_failure() {
		return results.Error(data, "用户名或密码错误")
	}

	the_user := rs.Data["user"].(model.User)

	rs = ctl.handler_user.Update_user(ctx, the_user.User_uuid, "is_login", true)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_user = rs.Data["user"].(model.User)

	request := iris.Map{
		"user_name":     user_name,
		"user_password": password,
	}

	rs = ctl.helper_token_granters.Grant(ctx, helper.TOKEN_GRANT_type_username_and_password, client_uuid, request)

	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	data["user"] = the_user
	data["access_token"] = rs.Data["access_token"]
	data["refresh_token"] = rs.Data["refresh_token"]

	return results.Ok(data)
}

func (ctl *User_http_controller) Get_all_users(ctx context.Context) *results.Result {

	data := results.API()

	rs := ctl.handler_user.Get_all_users(ctx)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	data["users"] = rs.Data["users"]
	return results.Ok(data)
}

func (ctl *User_http_controller) Logout(ctx context.Context, token_value, user_uuid string) *results.Result {

	data := results.API()

	rs := ctl.handler_user.Get_user_by_uuid(ctx, user_uuid)
	if rs.Is_failure() {
		tools.Log(rs.Msg)
		return results.Error(data, rs.Msg)
	}

	the_user := rs.Data["user"].(model.User)

	if the_user.Is_login == false {
		return results.Ok(data)
	}

	rs = ctl.handler_token_store.Read_access_token(ctx, token_value)
	if rs.Is_failure() {
		tools.Log(rs.Msg)
		return results.Error(data, rs.Msg)
	}

	the_access_token := rs.Data["token"].(model.Token)

	rs = ctl.handler_oauth.Get_oauth_by_uuid(ctx, the_access_token.Manager_oauth_uuid)
	if rs.Is_failure() {
		tools.Log(rs.Msg)
		return results.Error(data, rs.Msg)
	}

	the_oauth := rs.Data["oauth"].(model.Oauth)

	if the_oauth.Manager_user_uuid != user_uuid {
		return results.Error(data, results.Err_invalid_token_request.Error())
	}

	rs = ctl.handler_token_store.Remove_access_token(ctx, the_access_token.Token_value)
	if rs.Is_failure() {
		tools.Log(rs.Msg)
		return results.Error(data, rs.Msg)
	}

	rs = ctl.handler_token_store.Get_refresh_token_by_oauth(ctx, the_oauth.Oauth_uuid)
	if rs.Is_success() {
		the_refresh_token := rs.Data["token"].(model.Token)
		rs = ctl.handler_token_store.Remove_refresh_token(ctx, the_refresh_token.Token_value)
		if rs.Is_failure() {
			tools.Log(rs.Msg)
			return results.Error(data, rs.Msg)
		}
	}

	rs = ctl.handler_user.Update_user(ctx, the_user.User_uuid, "is_login", false)
	if rs.Is_failure() {
		tools.Log(rs.Msg)
		return results.Error(data, rs.Msg)
	}

	//tools.Log(data)
	return results.Ok(data)
}
