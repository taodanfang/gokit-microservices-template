package app

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 1.1) 定义 application 接口

type I_grpc_application interface {
	Make_route(string, endpoint.Endpoint)

	Get_db() *mongo.Database
	Get_router() *Router
	Get_logger() *log.Logger

	// gRPC payload 转码
	Decode_request(ctx context.Context, r interface{}) (request interface{}, err error)
	Encode_response(ctx context.Context, w interface{}) (response interface{}, err error)
}

type Router struct {
	Table map[string]grpc.Handler
}
