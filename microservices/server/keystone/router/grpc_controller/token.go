package grpc_controller

import (
	"cctable/common/libs/micro/server/grpc_server/app"
	"cctable/common/results"
	"cctable/microservices/server/keystone/business/model"
	"cctable/microservices/server/keystone/business/service"

	"context"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Token_grpc_controller struct {
	server app.I_grpc_application

	handler_token_store service.I_token_store
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_token_grpc_controller(server app.I_grpc_application) *Token_grpc_controller {
	ctl := &Token_grpc_controller{server: server}

	ctl.handler_token_store = service.New_token_store(server.Get_db())
	return ctl
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (ctl *Token_grpc_controller) Check_access_token(ctx context.Context, token_value string) *results.Result {

	data := results.API()

	rs := ctl.handler_token_store.Read_access_token(ctx, token_value)
	if rs.Is_failure() {
		return results.Error(data, results.Err_token_does_not_exist.Error())
	}

	the_token := rs.Data["token"].(model.Token)

	if the_token.Is_expired() {
		return results.Error(data, results.Err_token_has_expired.Error())
	}

	data["token"] = the_token
	return results.Ok(data)
}
