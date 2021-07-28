package http_controller

import (
	"cctable/common/libs/micro/server/http_server/app"
	"cctable/common/results"
	"cctable/microservices/server/keystone/business/helper"
	"cctable/microservices/server/keystone/business/model"
	"cctable/microservices/server/keystone/business/service"
	"context"
	"github.com/kataras/iris/v12"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Oauth_http_controller struct {
	server app.I_http_application

	handler_user          service.I_user_service
	handler_client        service.I_client_service
	handler_oauth         service.I_oauth_service
	handler_token_store   service.I_token_store
	helper_token_enhance  helper.I_token_enhance_helper
	helper_token          helper.I_token_helper
	helper_token_granters helper.I_token_grant_helper
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_oauth_http_controller(server app.I_http_application) *Oauth_http_controller {

	ctl := &Oauth_http_controller{server: server}

	ctl.handler_user = service.New_user_service(server.Get_db())
	ctl.handler_client = service.New_client_service(server.Get_db())
	ctl.handler_oauth = service.New_oauth_service(server.Get_db())
	ctl.handler_token_store = service.New_token_store(server.Get_db())
	ctl.helper_token_enhance = helper.New_jwt_token_enhancer(server.Get_db(), "jwt_secret")
	ctl.helper_token = helper.New_token_helper(server.Get_db(), ctl.helper_token_enhance)

	token_granter_with_user_and_password := helper.New_token_granter_with_username_and_password(server.Get_db(), ctl.helper_token)
	token_granter_with_refresh_token := helper.New_token_granter_with_refresh_token(server.Get_db(), ctl.helper_token)

	ctl.helper_token_granters = helper.New_token_granters(map[string]helper.I_token_grant_helper{
		helper.TOKEN_GRANT_type_username_and_password: token_granter_with_user_and_password,
		helper.TOKEN_GRANT_type_refresh_token:         token_granter_with_refresh_token,
	})

	return ctl
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (ctl *Oauth_http_controller) Grant_token_with_user_and_password(ctx context.Context, client_uuid string, user_name, user_password string) *results.Result {

	data := results.API()

	rs := ctl.handler_client.Get_client_by_uuid(ctx, client_uuid)
	if rs.Is_failure() {
		return results.Error(data, results.Err_client_does_not_exist.Error())
	}

	rs = ctl.handler_user.Check_user_by_name_and_password(ctx, user_name, user_password)
	if rs.Is_failure() {
		return results.Error(data, "用户名或密码错误")
	}

	request := iris.Map{
		"user_name":     user_name,
		"user_password": user_password,
	}

	rs = ctl.helper_token_granters.Grant(ctx, helper.TOKEN_GRANT_type_username_and_password, client_uuid, request)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	data["access_token"] = rs.Data["access_token"]
	data["refresh_token"] = rs.Data["refresh_token"]

	return results.Ok(data)
}

func (ctl *Oauth_http_controller) Grant_token_with_refresh_token(ctx context.Context, client_uuid string, token_value string) *results.Result {

	data := results.API()

	rs := ctl.handler_client.Get_client_by_uuid(ctx, client_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	rs = ctl.handler_token_store.Read_refresh_token(ctx, token_value)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	request := iris.Map{
		"refresh_token": token_value,
	}

	rs = ctl.helper_token_granters.Grant(ctx, helper.TOKEN_GRANT_type_refresh_token, client_uuid, request)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	data["access_token"] = rs.Data["access_token"]
	data["refresh_token"] = rs.Data["refresh_token"]

	return results.Ok(data)
}

func (ctl *Oauth_http_controller) Check_token(ctx context.Context, client_uuid string, access_token_value string) *results.Result {

	data := results.API()

	rs := ctl.handler_client.Get_client_by_uuid(ctx, client_uuid)
	if rs.Is_failure() {
		return results.Error(data, results.Err_client_does_not_exist.Error())
	}

	rs = ctl.handler_token_store.Read_access_token(ctx, access_token_value)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_token := rs.Data["token"].(model.Token)

	rs = ctl.handler_oauth.Get_oauth_by_uuid(ctx, the_token.Manager_oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_oauth := rs.Data["oauth"].(model.Oauth)

	rs = ctl.handler_client.Get_client_by_uuid(ctx, the_oauth.Manager_client_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_client := rs.Data["client"]

	rs = ctl.handler_user.Get_user_by_uuid(ctx, the_oauth.Manager_user_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_user := rs.Data["user"]

	data["client"] = the_client
	data["user"] = the_user
	return results.Ok(data)
}
