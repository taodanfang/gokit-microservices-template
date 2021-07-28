package controller

import (
	"cctable/common/libs/micro/client/grpc_client"
	"cctable/common/results"
	"cctable/common/tools"
	"cctable/config"
	"cctable/microservices/client/keystone/client_grpc/endpoints"
	"context"
	"google.golang.org/grpc"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Keystone_token_grpc_client_controller struct {
	grpc_conn   *grpc.ClientConn
	grpc_client *grpc_client.Transport_for_grpc_client
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_keystone_token_grpc_client_controller() (controller *Keystone_token_grpc_client_controller, err error) {
	config.Init_config()
	cfg := config.Get_config()

	conn, err := grpc.Dial(cfg.Keystone.Host+":"+cfg.Keystone.Grpc_port, grpc.WithInsecure())
	if err != nil {
		tools.Log(err)
		return nil, results.Err_connect_to_rpc_server_is_failure
	}

	tps := grpc_client.Make_grpc_client(conn)

	c := &Keystone_token_grpc_client_controller{
		grpc_conn:   conn,
		grpc_client: tps,
	}

	return c, nil
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (ctl *Keystone_token_grpc_client_controller) Check_access_token(ctx context.Context, token_value string) *grpc_client.Json_rpc_response {
	manager := grpc_client.Make_grpc_endpoint_client_manager(ctl.grpc_client)
	handler := endpoints.New_token_service_grpc_client(manager)
	rs := handler.Check_access_token(ctx, token_value)
	return rs
}

func (ctl *Keystone_token_grpc_client_controller) Close_client() {
	_ = ctl.grpc_conn.Close()
}
