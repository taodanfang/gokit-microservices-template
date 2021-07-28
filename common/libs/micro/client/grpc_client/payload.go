package grpc_client

import "github.com/kataras/iris/v12"

// --------------------------------------------------------------------
// Payload 定义
// --------------------------------------------------------------------

type Json_rpc_request struct {
	Method string   `json:"method"`
	Params iris.Map `json:"params"`
}

type body struct {
	Data    map[string]interface{} `json:"data"`
	Msg     string                 `json:"msg"`
	Expired bool                   `json:"expired"`
}

type Json_rpc_response struct {
	Success  bool   `json:"success"`
	Code     string `json:"code"`
	Response body   `json:"response"`
	Error    body   `json:"error"`
}
