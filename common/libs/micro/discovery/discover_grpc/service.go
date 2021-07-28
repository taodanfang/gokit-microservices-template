package discover_grpc

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Gprc_discovery_service struct {
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *Gprc_discovery_service) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *Gprc_discovery_service) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}
