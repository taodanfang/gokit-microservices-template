package endpoints

import (
	"cctable/common/libs/micro/server/grpc_server/app"
	"github.com/go-kit/kit/endpoint"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_grpc_endpoint_manager interface {
	Register_grpc_endpoint(route_name string, edp endpoint.Endpoint)
}

// --------------------------------------------------------------------
// 初始化方法
// --------------------------------------------------------------------

// 3.1) 初始化服务 controller

func Init_grpc_endpoint_managers(tps app.I_grpc_application) {
	New_token_grpc_manager(tps)
}
