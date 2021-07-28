package endpoints

import (
	"cctable/api"
	"cctable/common/libs/micro/server/http_server/app"
	"cctable/common/results"
	"cctable/microservices/server/device/router/http_controller"
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kataras/iris/v12"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Device_http_manager struct {
	transport_server app.I_http_application
	endpoints        map[string]endpoint.Endpoint

	ctl_device *http_controller.Device_http_controller
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_device_http_manager(tps app.I_http_application) {
	m := &Device_http_manager{transport_server: tps}
	m.ctl_device = http_controller.New_device_http_controller(tps)

	m.endpoints = make(map[string]endpoint.Endpoint, 0)

	m.Register_http_endpoint(api.EDP_http_device__device__register_device, m.register_device())
	m.Register_http_endpoint(api.EDP_http_device__device__get_all_devices, m.get_all_devices())
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (m *Device_http_manager) Register_http_endpoint(route_name string, edp endpoint.Endpoint) {
	tps := m.transport_server
	m.endpoints[route_name] = edp
	tps.Make_route(route_name, m.endpoints[route_name])
}

// --------------------------------------------------------------------
// Endpoint 方法
// --------------------------------------------------------------------

func (m *Device_http_manager) register_device() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		data := results.API()

		err = Check_token_authorization(ctx, m.transport_server)
		if err != nil {
			rs := results.Error(data, err.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		device_name := ""
		device_code := ""

		rq := request.(app.Http_payload_request)
		if value, ok := rq.Params["device_name"]; ok {
			device_name = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		if value, ok := rq.Params["device_code"]; ok {
			device_code = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		rs := m.ctl_device.Register_device(ctx, device_name, device_code)
		if rs.Is_failure() {
			rs = results.Error(data, rs.Msg)
			return app.Http_payload_response{Result: *rs}, nil
		}

		data["result"] = iris.Map{
			"device": rs.Data["device"],
		}
		return app.Http_payload_response{Result: *results.Ok(data)}, nil
	}
}

func (m *Device_http_manager) get_all_devices() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		data := results.API()

		err = Check_token_authorization(ctx, m.transport_server)
		if err != nil {
			rs := results.Error(data, err.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		rs := m.ctl_device.Get_all_devices(ctx)
		if rs.Is_failure() {
			rs = results.Error(data, rs.Msg)
			return app.Http_payload_response{Result: *rs}, nil
		}

		data["result"] = rs.Data
		return app.Http_payload_response{Result: *results.Ok(data)}, nil
	}
}
