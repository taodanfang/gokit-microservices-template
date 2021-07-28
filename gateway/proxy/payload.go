package proxy

import (
	"bytes"
	"cctable/common/results"
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"net/http"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 1.2) 定义 payload 结构

type Http_gateway_request struct {
	Authorization string
	Params        iris.Map `json:"params"`
}

type Http_gateway_response struct {
	Result results.Result `json:"result"`
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func encode_http_request(_ context.Context, r *http.Request, request interface{}) error {

	//pp.Println("encode_request: ", request)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request.(Http_gateway_request)); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	r.Header.Set("Authorization", request.(Http_gateway_request).Authorization)

	return nil
}

func encode_http_response(_ context.Context, w http.ResponseWriter, response interface{}) error {

	err := json.NewEncoder(w).Encode(response)

	//pp.Println("encode_response: ", response)

	return err
}

func decode_http_request_FuncBuilder(request Http_gateway_request) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (_request interface{}, err error) {

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, err
		}

		request.Authorization = r.Header.Get("Authorization")

		//pp.Println("builder_decode_request: ", request)

		return request, nil
	}
}

func decode_http_response_FuncBuilder(response interface{}) httptransport.DecodeResponseFunc {
	return func(_ context.Context, r *http.Response) (_response interface{}, err error) {
		//pp.Println("decode_response: ", r)

		if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
			return nil, err
		}

		//pp.Println("builder_decode_response: ", response)
		return response, nil
	}
}
