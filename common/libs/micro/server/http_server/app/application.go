package app

import (
	"cctable/common/results"
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 1.1) 定义 application 接口

type I_http_application interface {
	Make_route(string, endpoint.Endpoint)

	Get_db() *mongo.Database
	Get_router() *iris.Application

	Decode_request(ctx context.Context, request *http.Request) (interface{}, error)
	Encode_response(ctx context.Context, w http.ResponseWriter, response interface{}) error
	Encode_error(ctx context.Context, err error, w http.ResponseWriter)
}

// 1.2) 定义 payload 结构

type Http_payload_request struct {
	Params iris.Map `json:"params`
}

type Http_payload_response struct {
	Result results.Result `json:"result"`
}
