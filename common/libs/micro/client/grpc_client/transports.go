package grpc_client

import (
	"cctable/common/libs/micro/pb"
	"cctable/common/results"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	go_kit_grpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_grpc_client interface {
	Make_endpoint(service_name, method string) endpoint.Endpoint

	// gRPC payload 转码
	Encode_request(ctx context.Context, w interface{}) (request interface{}, err error)
	Decode_response(ctx context.Context, r interface{}) (response interface{}, err error)
}

type Transport_for_grpc_client struct {
	conn *grpc.ClientConn
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

// 2.2) 创建 application

func Make_grpc_client(conn *grpc.ClientConn) *Transport_for_grpc_client {

	tps := &Transport_for_grpc_client{
		conn: conn,
	}

	return tps
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (t *Transport_for_grpc_client) Make_endpoint(service_name, method string) endpoint.Endpoint {
	client := go_kit_grpc.NewClient(
		t.conn,
		service_name,
		method,
		t.Encode_request,
		t.Decode_response,
		pb.JsonResponse{},
	)

	return client.Endpoint()
}

func (t *Transport_for_grpc_client) Encode_request(ctx context.Context, w interface{}) (request interface{}, err error) {

	//pp.Println("encode_request: ", w)

	rq := w.(*Json_rpc_request)

	//tools.Log(rq)

	if rq.Method == "" {
		return nil, results.Err_invalid_method_with_rpc_request
	}

	byte_request, _ := json.Marshal(rq)

	return &pb.JsonRequest{Params: byte_request}, nil
}

func (t *Transport_for_grpc_client) Decode_response(ctx context.Context, r interface{}) (response interface{}, err error) {

	//pp.Println("decode_response: ", r)

	req := r.(*pb.JsonResponse)

	//tools.Log(req)

	var api_result Json_rpc_response
	err = json.Unmarshal(req.Result, &api_result)
	if err != nil {
		return nil, err
	}

	//pp.Println(api_result)

	return &api_result, nil
}
