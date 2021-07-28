package http_controller

import (
	"cctable/common/libs/micro/server/http_server/app"
	"cctable/common/results"
	"cctable/microservices/server/device/business/model"
	"cctable/microservices/server/device/business/service"
	"context"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Device_http_controller struct {
	server app.I_http_application

	handler_device service.I_device_service
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_device_http_controller(server app.I_http_application) *Device_http_controller {
	ctl := &Device_http_controller{server: server}
	ctl.handler_device = service.New_device_service(server.Get_db())
	return ctl
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (ctl *Device_http_controller) Register_device(ctx context.Context, device_name, device_code string) *results.Result {

	data := results.API()

	rs := ctl.handler_device.Check_device_by_name_and_code(ctx, device_name, device_code)
	if rs.Is_success() {
		return results.Ok(data)
	}

	rs = ctl.handler_device.New_device(ctx, device_name, device_code)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_device := rs.Data["device"].(model.Device)

	data["device"] = the_device

	return results.Ok(data)
}

func (ctl *Device_http_controller) Get_all_devices(ctx context.Context) *results.Result {

	data := results.API()

	rs := ctl.handler_device.Get_all_devices(ctx)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	data["devices"] = rs.Data["devices"]
	return results.Ok(data)
}
