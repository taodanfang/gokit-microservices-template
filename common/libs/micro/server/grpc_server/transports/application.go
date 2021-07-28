package transports

import (
	"cctable/common/libs/micro/pb"
	"cctable/common/libs/micro/server/grpc_server/app"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 2.1) 实现 application 接口

type Transport_for_grpc_application struct {
	db             *mongo.Database
	logger         log.Logger
	router         *app.Router
	server_options []grpc.ServerOption
}

type body struct {
	Data    map[string]interface{} `json:"data"`
	Msg     string                 `json:"msg"`
	Expired bool                   `json:"expired"`
}

type json_response struct {
	Success  bool   `json:"success"`
	Code     string `json:"code"`
	Response body   `json:"response"`
	Error    body   `json:"error"`
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

// 2.2) 创建 application

func Make_grpc_application(db *mongo.Database, logger log.Logger) *Transport_for_grpc_application {
	router := app.Router{
		Table: make(map[string]grpc.Handler, 0),
	}
	return &Transport_for_grpc_application{
		db:     db,
		logger: logger,
		router: &router,
	}
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (t *Transport_for_grpc_application) Make_route(method string, endpoint endpoint.Endpoint) {
	handler := grpc.NewServer(
		endpoint,
		t.Decode_request,
		t.Encode_response,
		t.server_options...,
	)

	t.router.Table[method] = handler

	//pp.Println(t.router.Table)
}

func (t *Transport_for_grpc_application) Get_db() *mongo.Database {
	return t.db
}

func (t *Transport_for_grpc_application) Get_logger() *log.Logger {
	return &t.logger
}

func (t *Transport_for_grpc_application) Get_router() *app.Router {
	return t.router
}

func (t *Transport_for_grpc_application) Decode_request(ctx context.Context, r interface{}) (request interface{}, err error) {
	req := r.(*pb.JsonRequest)

	var api_request iris.Map
	err = json.Unmarshal(req.Params, &api_request)
	if err != nil {
		return nil, err
	}

	method, ok := api_request["method"]
	if !ok {
		return nil, err
	}

	params, ok := api_request["params"]
	if !ok {
		params = iris.Map{}
	}

	return &app.Grpc_payload_request{Method: method.(string), Params: params.(iris.Map)}, nil
}

func (t *Transport_for_grpc_application) Encode_response(ctx context.Context, w interface{}) (response interface{}, err error) {

	rsp := w.(app.Grpc_payload_response)

	var api_response interface{}

	if rsp.Result.Is_success() {
		api_response = make_json_success_response(w)
	} else {
		api_response = make_json_failure_response(w)
	}

	byte_result, _ := json.Marshal(api_response)
	return &pb.JsonResponse{Result: byte_result}, nil
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func make_json_success_response(result interface{}) interface{} {

	rs := result.(app.Grpc_payload_response).Result
	return &json_response{
		Success: true,
		Code:    "server-api-ok",
		Response: body{
			Data:    rs.Data,
			Msg:     rs.Msg,
			Expired: rs.Api_is_expired(),
		},
		Error: body{
			Data:    iris.Map{},
			Msg:     "",
			Expired: false,
		},
	}
}

func make_json_failure_response(result interface{}) interface{} {

	rs := result.(app.Grpc_payload_response).Result
	return &json_response{
		Success: false,
		Code:    "server-api-error",
		Response: body{
			Data:    iris.Map{},
			Msg:     "",
			Expired: false,
		},
		Error: body{
			Data:    rs.Data,
			Msg:     rs.Msg,
			Expired: rs.Api_is_expired(),
		},
	}
}
