package pb_server

import (
	"cctable/api"
	"cctable/common/libs/micro/pb"
	"cctable/common/libs/micro/server/grpc_server/app"
	"context"
	"github.com/go-kit/kit/transport/grpc"
)

// --------------------------------------------------------------------
// 对象定义
// --------------------------------------------------------------------

// 3.1) 实现 pb.server

type Token_service_grpc_router struct {
	transport_server   app.I_grpc_application
	check_access_token grpc.Handler
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

// 3.2) 创建 pb.server

func New_token_service_grpc_router(tps app.I_grpc_application) *Token_service_grpc_router {
	router := Token_service_grpc_router{transport_server: tps}
	router.init_router()
	return &router
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

// 3.3) 将 endpoint 绑定到 application 的 handler

func (s *Token_service_grpc_router) init_router() {

	router_table := s.transport_server.Get_router().Table

	for method, handler := range router_table {
		switch method {
		case api.EDP_grpc_keystone_token_check_access_token:
			s.check_access_token = handler
		}
	}
}

// --------------------------------------------------------------------
// gRPC server 调用方法（实现 pb 中 xxxServiceServer interface）
// --------------------------------------------------------------------

func (s *Token_service_grpc_router) CheckAccessToken(ctx context.Context, r *pb.JsonRequest) (*pb.JsonResponse, error) {
	//pp.Println("string_server.Concat", r)

	_, resp, err := s.check_access_token.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.JsonResponse), nil
}
