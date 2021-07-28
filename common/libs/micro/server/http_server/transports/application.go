package transports

import (
	"cctable/common/libs/micro/server/http_server/app"
	"cctable/common/results"
	"cctable/common/tools"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	go_kit_http "github.com/go-kit/kit/transport/http"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

// --------------------------------------------------------------------
// 常量定义
// --------------------------------------------------------------------

const (
	Authorization_value_key = "Authorization_value_key"
	Authorization_error_key = "Authorization_error_key"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 2.1) 实现 application 接口

type Transport_for_http_application struct {
	DB             *mongo.Database
	Router         *iris.Application
	Server_options []go_kit_http.ServerOption
}

type json_response struct {
	Success  bool     `json:"success"`
	Code     string   `json:"code"`
	Response iris.Map `json:"response"`
	Error    iris.Map `json:"error"`
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

// 2.2) 创建 application

func Make_http_application(router *iris.Application, db *mongo.Database) *Transport_for_http_application {

	tps := &Transport_for_http_application{Router: router, DB: db}

	client_authorization_options := []go_kit_http.ServerOption{
		go_kit_http.ServerBefore(make_authorization_context()),
		go_kit_http.ServerErrorHandler(transport.NewLogErrorHandler(tools.Get_gokit_logger())),
		go_kit_http.ServerErrorEncoder(tps.Encode_error),
	}

	tps.Server_options = client_authorization_options

	return tps
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (t *Transport_for_http_application) Make_route(route_name string, endpoint endpoint.Endpoint) {

	t.Router.Post(route_name, iris.FromStd(
		go_kit_http.NewServer(
			endpoint,
			t.Decode_request,
			t.Encode_response,
			t.Server_options...,
		)))

}

func (t *Transport_for_http_application) Get_db() *mongo.Database {
	return t.DB
}

func (t *Transport_for_http_application) Get_router() *iris.Application {
	return t.Router
}

func (t *Transport_for_http_application) Decode_request(ctx context.Context, request *http.Request) (interface{}, error) {

	//pp.Println("request_header: ", request.Header)

	var api_request app.Http_payload_request
	err := json.NewDecoder(request.Body).Decode(&api_request)
	if err != nil {
		api_request = app.Http_payload_request{}
	}

	//tools.Log("decode_request: ", api_request)

	return api_request, nil
}

func (t *Transport_for_http_application) Encode_response(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	var iris_rsp interface{}

	rsp := response.(app.Http_payload_response)

	if rsp.Result.Is_success() {
		iris_rsp = make_json_success_response(response)
	} else {
		iris_rsp = make_json_failure_response(response)
	}
	return json.NewEncoder(w).Encode(iris_rsp)
}

func (t *Transport_for_http_application) Encode_error(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(iris.Map{
		"error": err.Error(),
	})
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func make_json_success_response(result interface{}) interface{} {

	rs := result.(app.Http_payload_response).Result
	code := "server-api-ok"
	if rs.Api_is_expired() {
		code = "server-api-expired"
	}
	return &json_response{
		Success: true,
		Code:    code,
		Response: iris.Map{
			"data": rs.Data,
			"msg":  rs.Msg,
		},
		Error: iris.Map{},
	}
}

func make_json_failure_response(result interface{}) interface{} {

	rs := result.(app.Http_payload_response).Result
	code := "server-api-error"
	if rs.Api_is_expired() {
		code = "server-api-expired"
	}
	return &json_response{
		Success:  false,
		Code:     code,
		Response: iris.Map{},
		Error: iris.Map{
			"data": rs.Data,
			"msg":  rs.Msg,
		},
	}
}

func make_authorization_context() go_kit_http.RequestFunc {
	return func(ctx context.Context, rq *http.Request) context.Context {

		//tools.Log("decode_request_header: ", rq.Header)

		authorization_string := rq.Header.Get("Authorization")
		authorization_string = strings.TrimSpace(authorization_string)

		if authorization_string == "" {
			return context.WithValue(ctx, Authorization_error_key, results.Err_lost_authorization_with_request)
		}

		authorization_string_split := strings.Split(authorization_string, " ")
		if len(authorization_string_split) != 2 {
			return context.WithValue(ctx, Authorization_error_key, results.Err_invalid_authorization_header)
		}

		if authorization_string_split[0] == "Basic" {
			client_uuid, client_secret, ok := rq.BasicAuth()

			if ok {
				auth := iris.Map{
					"auth_type":     "oauth",
					"client_uuid":   client_uuid,
					"client_secret": client_secret,
				}

				//tools.Log("auth:", auth)
				return context.WithValue(ctx, Authorization_value_key, auth)
			}

			return context.WithValue(ctx, Authorization_error_key, results.Err_invalid_client_request)
		}

		if authorization_string_split[0] == "Bearer" {
			auth := iris.Map{
				"auth_type":   "token",
				"token_value": authorization_string_split[1],
			}

			//tools.Log("auth:", auth)

			return context.WithValue(ctx, Authorization_value_key, auth)
		}

		return context.WithValue(ctx, Authorization_error_key, results.Err_invalid_authorization_header)
	}
}
