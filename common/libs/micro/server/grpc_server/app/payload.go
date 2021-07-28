package app

import (
	"cctable/common/results"
	"github.com/kataras/iris/v12"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 1.2) 定义 payload 结构

type Grpc_payload_request struct {
	Method string   `json:"method"` // 全部小写
	Params iris.Map `json:"params"`
}

type Grpc_payload_response struct {
	Result results.Result `json:"result"`
}

// --------------------------------------------------------------------
// 导出方法
// --------------------------------------------------------------------

func (rq Grpc_payload_request) Check_method(method_name string) (err error) {

	if rq.Method != method_name {

		return results.Err_invalid_method_with_rpc_request
	}

	return nil
}

func (rq Grpc_payload_request) Check_param(param_name string) (err error) {

	_, ok := rq.Params[param_name]
	if !ok {
		return results.Err_invalid_params_with_rpc_request
	}

	return nil
}

func (rq Grpc_payload_request) Check_request(method_name string, param_names ...string) error {
	if rq.Method != method_name {
		return results.Err_invalid_method_with_rpc_request
	}

	for _, param_name := range param_names {
		_, ok := rq.Params[param_name]
		if !ok {
			return results.Err_invalid_params_with_rpc_request
		}
	}

	return nil
}
