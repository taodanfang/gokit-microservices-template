package helper

import (
	"cctable/common/results"
	"cctable/microservices/server/keystone/business/model"
	"cctable/microservices/server/keystone/business/service"
	"context"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 常量定义
// --------------------------------------------------------------------

const (
	TOKEN_GRANT_type_username_and_password = "password"
	TOKEN_GRANT_type_refresh_token         = "refresh_token"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_token_grant_helper interface {
	Grant(ctx context.Context, grant_type string, client_uuid string, request iris.Map) *results.Result
}

type Token_granters struct {
	Token_grant_dict map[string]I_token_grant_helper
}

type Token_granter_with_username_and_password struct {
	support_grant_type     string
	handler_user_service   service.I_user_service
	handler_client_service service.I_client_service
	handler_oauth_service  service.I_oauth_service
	helper_token_service   I_token_helper
}

type Token_granter_with_refresh_token struct {
	support_grant_type   string
	handler_token_store  service.I_token_store
	helper_token_service I_token_helper
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_token_granters(token_grant_dict map[string]I_token_grant_helper) *Token_granters {
	return &Token_granters{
		Token_grant_dict: token_grant_dict,
	}
}

func New_token_granter_with_username_and_password(db *mongo.Database, token_service I_token_helper) *Token_granter_with_username_and_password {
	return &Token_granter_with_username_and_password{
		support_grant_type:     TOKEN_GRANT_type_username_and_password,
		handler_user_service:   service.New_user_service(db),
		handler_client_service: service.New_client_service(db),
		handler_oauth_service:  service.New_oauth_service(db),
		helper_token_service:   token_service,
	}
}

func New_token_granter_with_refresh_token(db *mongo.Database, token_service I_token_helper) *Token_granter_with_refresh_token {
	return &Token_granter_with_refresh_token{
		support_grant_type:   TOKEN_GRANT_type_refresh_token,
		handler_token_store:  service.New_token_store(db),
		helper_token_service: token_service,
	}
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (hlp *Token_granters) Grant(ctx context.Context, grant_type string, client_uuid string, request iris.Map) *results.Result {
	data := results.API()

	dispatcher_granter := hlp.Token_grant_dict[grant_type]
	if dispatcher_granter == nil {
		return results.Error(data, results.Err_does_not_support_the_grant_type.Error())
	}

	return dispatcher_granter.Grant(ctx, grant_type, client_uuid, request)
}

func (hlp *Token_granter_with_username_and_password) Grant(ctx context.Context, grant_type string, client_uuid string, request iris.Map) *results.Result {

	data := results.API()

	if grant_type != hlp.support_grant_type {
		return results.Error(data, results.Err_does_not_support_the_grant_type.Error())
	}

	user_name := ""
	user_password := ""

	if value, ok := request["user_name"]; ok {
		user_name = value.(string)
	}

	if value, ok := request["user_password"]; ok {
		user_password = value.(string)
	}

	if user_name == "" || user_password == "" {
		return results.Error(data, results.Err_invalid_user_name_and_password.Error())
	}

	rs := hlp.handler_user_service.Check_user_by_name_and_password(ctx, user_name, user_password)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_user := rs.Data["user"].(model.User)

	rs = hlp.handler_oauth_service.Get_oauth_by_user_and_client_uuid(ctx, the_user.User_uuid, client_uuid)
	//pp.Println(rs)

	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_oauth := rs.Data["oauth"].(model.Oauth)
	//pp.Println(the_oauth)

	rs = hlp.helper_token_service.Create_access_token(ctx, the_oauth.Oauth_uuid)
	//pp.Println(rs)

	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	data["access_token"] = rs.Data["access_token"]
	data["refresh_token"] = rs.Data["refresh_token"]

	return results.Ok(data)
}

func (hlp *Token_granter_with_refresh_token) Grant(ctx context.Context, grant_type string, client_uuid string, request iris.Map) *results.Result {
	data := results.API()

	if grant_type != hlp.support_grant_type {
		return results.Error(data, results.Err_does_not_support_the_grant_type.Error())
	}

	refresh_token_value := ""
	if value, ok := request["refresh_token"]; ok {
		refresh_token_value = value.(string)
	}

	if refresh_token_value == "" {
		return results.Error(data, results.Err_invalid_refresh_token_value.Error())
	}

	rs := hlp.handler_token_store.Read_refresh_token(ctx, refresh_token_value)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_refresh_token := rs.Data["token"].(model.Token)

	return hlp.helper_token_service.Refresh_access_token(ctx, the_refresh_token.Manager_oauth_uuid)
}
