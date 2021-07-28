package discover_http

import (
	"cctable/common/tools"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	go_kit_http "github.com/go-kit/kit/transport/http"
	"github.com/kataras/iris/v12"
	"net/http"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------
type I_http_discovery interface {
	Make_route(route_name string, endpoint endpoint.Endpoint, decode_request_func go_kit_http.DecodeRequestFunc)

	Get_discovery_client() I_discovery_http_client
	Get_router() *iris.Application

	Encode_response(ctx context.Context, w http.ResponseWriter, response interface{}) error
	Encode_error(ctx context.Context, err error, w http.ResponseWriter)
}

type Transport_for_http_discovery struct {
	discovery_client I_discovery_http_client
	router           *iris.Application
	server_options   []go_kit_http.ServerOption
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func Make_http_discovery(router *iris.Application, discovery_client I_discovery_http_client) {
	tps := &Transport_for_http_discovery{router: router, discovery_client: discovery_client}

	server_options := []go_kit_http.ServerOption{
		go_kit_http.ServerErrorHandler(transport.NewLogErrorHandler(tools.Get_gokit_logger())),
		go_kit_http.ServerErrorEncoder(tps.Encode_error),
	}

	tps.server_options = server_options

	New_http_discovery_controller(tps)
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (t *Transport_for_http_discovery) Make_route(route_name string, endpoint endpoint.Endpoint, decode_request_func go_kit_http.DecodeRequestFunc) {

	t.router.Get("/"+route_name, iris.FromStd(
		go_kit_http.NewServer(
			endpoint,
			decode_request_func,
			t.Encode_response,
			t.server_options...,
		)))

}

func (t *Transport_for_http_discovery) Get_discovery_client() I_discovery_http_client {
	return t.discovery_client
}

func (t *Transport_for_http_discovery) Get_router() *iris.Application {
	return t.router
}

func (t *Transport_for_http_discovery) Encode_response(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func (t *Transport_for_http_discovery) Encode_error(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(iris.Map{
		"error": err.Error(),
	})
}
