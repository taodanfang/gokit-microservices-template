package helper

import (
	"cctable/common/results"
	"cctable/common/tools"
	"cctable/microservices/server/keystone/business/model"
	"cctable/microservices/server/keystone/business/service"
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 6.1) 定义组合服务接口

type I_token_helper interface {
	Create_access_token(ctx context.Context, oauth_uuid string) *results.Result
	Refresh_access_token(ctx context.Context, oauth_uuid string) *results.Result
}

// 6.2) 实现组合服务接口

type Token_helper struct {
	handler_token_store    service.I_token_store
	helper_token_enhance   I_token_enhance_helper
	handler_user_service   service.I_user_service
	handler_client_service service.I_client_service
	handler_oauth_service  service.I_oauth_service
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

// 6.3) 创建组合服务

func New_token_helper(
	db *mongo.Database,
	token_enhance I_token_enhance_helper) *Token_helper {

	return &Token_helper{
		handler_token_store:    service.New_token_store(db),
		helper_token_enhance:   token_enhance,
		handler_user_service:   service.New_user_service(db),
		handler_client_service: service.New_client_service(db),
		handler_oauth_service:  service.New_oauth_service(db),
	}
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (hlp *Token_helper) Create_access_token(ctx context.Context, oauth_uuid string) *results.Result {
	data := results.API()

	rs := hlp.handler_oauth_service.Get_oauth_by_uuid(ctx, oauth_uuid)
	//pp.Println(rs)

	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	var has_access_token = false
	var the_access_token model.Token

	rs = hlp.handler_token_store.Get_access_token_by_oauth(ctx, oauth_uuid)
	//pp.Println(rs)

	if rs.Is_success() {
		has_access_token = true
		the_access_token = rs.Data["token"].(model.Token)
	}

	if has_access_token {
		if the_access_token.Is_expired() == false {
			data["access_token"] = the_access_token
			rs = hlp.handler_token_store.Get_refresh_token_by_oauth(ctx, oauth_uuid)
			if rs.Is_failure() {
				return results.Error(data, rs.Msg)
			}
			data["refresh_token"] = rs.Data["token"]
			return results.Ok(data)
		}

		rs = hlp.handler_token_store.Remove_access_token(ctx, the_access_token.Token_value)
		if rs.Is_failure() {
			return results.Error(data, rs.Msg)
		}

		rs = hlp.handler_token_store.Get_refresh_token_by_oauth(ctx, oauth_uuid)
		if rs.Is_success() {
			the_refresh_token := rs.Data["token"].(model.Token)
			rs = hlp.handler_token_store.Remove_refresh_token(ctx, the_refresh_token.Token_value)
		}
	}

	rs = hlp.handler_token_store.Get_refresh_token_by_oauth(ctx, oauth_uuid)
	if rs.Is_success() {
		the_refresh_token := rs.Data["token"].(model.Token)
		rs = hlp.handler_token_store.Remove_refresh_token(ctx, the_refresh_token.Token_value)
	}

	rs = hlp.create_refresh_token(ctx, oauth_uuid)
	//pp.Println(rs)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_new_refresh_token := rs.Data["token"].(model.Token)

	rs = hlp.create_access_token(ctx, oauth_uuid)
	//pp.Println(rs)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_new_access_token := rs.Data["token"].(model.Token)

	data["access_token"] = the_new_access_token
	data["refresh_token"] = the_new_refresh_token

	//pp.Println(data)

	return results.Ok(data)
}

func (hlp *Token_helper) Refresh_access_token(ctx context.Context, oauth_uuid string) *results.Result {

	data := results.API()

	var the_refresh_token model.Token

	rs := hlp.handler_token_store.Get_refresh_token_by_oauth(ctx, oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_refresh_token = rs.Data["token"].(model.Token)

	if the_refresh_token.Is_expired() {

		rs = hlp.handler_token_store.Remove_refresh_token(ctx, the_refresh_token.Token_value)
		if rs.Is_failure() {
			return results.Error(data, rs.Msg)
		}

		rs = hlp.handler_token_store.Get_access_token_by_oauth(ctx, oauth_uuid)
		if rs.Is_failure() {
			return results.Error(data, rs.Msg)
		}

		the_access_token := rs.Data["token"].(model.Token)

		rs = hlp.handler_token_store.Remove_access_token(ctx, the_access_token.Token_value)
		if rs.Is_failure() {
			return results.Error(data, rs.Msg)
		}

		return results.Error(data, results.Err_token_has_expired.Error())
	}

	rs = hlp.handler_token_store.Remove_refresh_token(ctx, the_refresh_token.Token_value)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	rs = hlp.handler_token_store.Get_access_token_by_oauth(ctx, oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_access_token := rs.Data["token"].(model.Token)

	rs = hlp.handler_token_store.Remove_access_token(ctx, the_access_token.Token_value)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	rs = hlp.create_refresh_token(ctx, oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_new_refresh_token := rs.Data["token"].(model.Token)

	rs = hlp.create_access_token(ctx, oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_new_access_token := rs.Data["token"].(model.Token)

	data["access_token"] = the_new_access_token
	data["refresh_token"] = the_new_refresh_token

	return results.Ok(data)
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (hlp *Token_helper) create_refresh_token(ctx context.Context, oauth_uuid string) *results.Result {

	data := results.API()

	//pp.Println(data, oauth_uuid)

	rs := hlp.handler_oauth_service.Get_oauth_by_uuid(ctx, oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_oauth := rs.Data["oauth"].(model.Oauth)

	//pp.Println(the_oauth)

	rs = hlp.handler_client_service.Get_client_by_uuid(ctx, the_oauth.Manager_client_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_client := rs.Data["client"].(model.Client)

	s, _ := time.ParseDuration(strconv.Itoa(the_client.Refresh_token_validity_seconds) + "s")
	the_token_expired_time := time.Now().Add(s)
	the_token_value := tools.NewUUID()

	rs = hlp.handler_token_store.Store_refresh_token(ctx, model.TOKEN_type_uuid, the_token_value, the_token_expired_time, oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_token := rs.Data["token"].(model.Token)

	if hlp.helper_token_enhance != nil {
		rs = hlp.helper_token_enhance.Enhance(ctx, the_token.Token_uuid)
		if rs.Is_failure() {
			return results.Error(data, rs.Msg)
		}

		the_token = rs.Data["token"].(model.Token)
		//pp.Println(the_token)
	}

	data["token"] = the_token
	return results.Ok(data)
}

func (hlp *Token_helper) create_access_token(ctx context.Context, oauth_uuid string) *results.Result {

	data := results.API()

	//pp.Println(data, oauth_uuid)

	rs := hlp.handler_oauth_service.Get_oauth_by_uuid(ctx, oauth_uuid)
	//pp.Println(rs)

	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_oauth := rs.Data["oauth"].(model.Oauth)

	rs = hlp.handler_client_service.Get_client_by_uuid(ctx, the_oauth.Manager_client_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_client := rs.Data["client"].(model.Client)

	s, _ := time.ParseDuration(strconv.Itoa(the_client.Access_token_validity_seconds) + "s")
	the_token_expired_time := time.Now().Add(s)
	the_token_value := tools.NewUUID()

	rs = hlp.handler_token_store.Store_access_token(ctx, model.TOKEN_type_uuid, the_token_value, the_token_expired_time, oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_token := rs.Data["token"].(model.Token)

	if hlp.helper_token_enhance != nil {
		rs = hlp.helper_token_enhance.Enhance(ctx, the_token.Token_uuid)
		if rs.Is_failure() {
			return results.Error(data, rs.Msg)
		}

		the_token = rs.Data["token"].(model.Token)
	}

	data["token"] = the_token
	return results.Ok(data)
}
