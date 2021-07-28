package endpoints

import (
	"cctable/common/libs/micro/server/http_server/app"
	"cctable/common/results"
	"cctable/microservices/server/message/business/service"
	"context"
	"github.com/kataras/iris/v12"

	"github.com/go-kit/kit/endpoint"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Room_http_manager struct {
	transport_server app.I_http_application
	endpoints        map[string]endpoint.Endpoint

	handler_room service.I_room_service
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_room_http_manager(tps app.I_http_application) {
	m := &Room_http_manager{transport_server: tps}

	m.handler_room = service.New_room_service(tps.Get_db())

	m.endpoints = make(map[string]endpoint.Endpoint, 0)

	m.register_endpoint("/room/check_room", m.check_room())
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (m *Room_http_manager) register_endpoint(route_name string, edp endpoint.Endpoint) {
	tps := m.transport_server
	m.endpoints[route_name] = edp
	tps.Make_route(route_name, m.endpoints[route_name])
}

// --------------------------------------------------------------------
// Endpoint 方法
// --------------------------------------------------------------------

// 3.5) 定义服务 endpoint

func (m *Room_http_manager) check_room() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		data := results.API()

		room_name := ""

		rq := request.(app.Http_payload_request)
		if value, ok := rq.Params["room_name"]; ok {
			room_name = value.(string)
		} else {
			rs := results.Error(data, results.Err_lost_parameters_with_request.Error())
			return app.Http_payload_response{Result: *rs}, nil
		}

		rs := m.handler_room.Check_room_by_name(ctx, room_name)
		if rs.Is_failure() {
			rs = results.Error(data, rs.Msg)
			return app.Http_payload_response{Result: *rs}, nil
		}

		data["result"] = iris.Map{
			"room": rs.Data["room"],
		}
		return app.Http_payload_response{Result: *results.Ok(data)}, nil
	}
}
