package endpoints

import (
	"cctable/common/libs/micro/server/http_server/app"
	"github.com/go-kit/kit/endpoint"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_http_endpoint_manager interface {
	Register_http_endpoint(route_name string, edp endpoint.Endpoint)
}

// --------------------------------------------------------------------
// 初始化方法
// --------------------------------------------------------------------

// 3.1) 初始化服务 controller
func Init_http_endpoint_managers(tps app.I_http_application) {
	New_room_http_manager(tps)
}
