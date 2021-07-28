package proxy

import (
	"cctable/common/results"
	"cctable/common/tools"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"net/url"
	"strings"
	"time"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Service_endpoint_manager struct {
	instance_manager *Service_instance_manager

	endpoint_map map[string]*endpoint_handler
}

type endpoint_handler struct {
	transport_type string
	service_name   string
	path           string
	method         string

	endpoint endpoint.Endpoint
	handler  *httptransport.Server
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_service_endpoint_manager(instance_manager *Service_instance_manager) *Service_endpoint_manager {

	return &Service_endpoint_manager{
		instance_manager: instance_manager,
		endpoint_map:     make(map[string]*endpoint_handler, 0),
	}
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *Service_endpoint_manager) Register_http_endpoint(service_name string, path string, method string) error {

	an_endpoint_handler := &endpoint_handler{
		transport_type: "http",
		service_name:   service_name,
		path:           path,
		method:         method,
	}

	the_instancer := s.instance_manager.Get_instance(service_name)

	//tools.Log("register_endpoint: ", the_instancer)

	if the_instancer == nil {
		return results.Err_service_does_not_exist
	}

	//pp.Println("instancer: ", the_instancer)

	the_endpointer := sd.NewEndpointer(
		the_instancer, http_service_factory_builder(path, method), tools.Get_gokit_logger())
	the_endpoint := lb.Retry(3, 3*time.Second, lb.NewRoundRobin(the_endpointer))

	an_endpoint_handler.endpoint = the_endpoint

	the_handler := httptransport.NewServer(
		the_endpoint,
		decode_http_request_FuncBuilder(Http_gateway_request{}),
		encode_http_response,
	)

	an_endpoint_handler.handler = the_handler

	s.endpoint_map[service_name+"@"+path] = an_endpoint_handler

	//pp.Println(s.endpoint_map)

	return nil
}

func (s *Service_endpoint_manager) Get_handler(service_name string, path string) (*httptransport.Server, error) {

	an_endpoint_handler, ok := s.endpoint_map[service_name+"@"+path]
	if !ok {
		return nil, results.Err_service_does_not_exist
	}

	return an_endpoint_handler.handler, nil
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func http_service_factory_builder(path string, method string) sd.Factory {
	return func(instance string) (e endpoint.Endpoint, closer io.Closer, err error) {

		http_prefix := "http://"
		if !strings.HasPrefix(instance, http_prefix) {
			instance = http_prefix + instance
		}
		//pp.Println("instance: ", instance)
		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}

		tgt.Path = path

		//pp.Println("tgt: ", tgt)
		return httptransport.NewClient(
			method, tgt,
			encode_http_request, decode_http_response_FuncBuilder(Http_gateway_response{}),
		).Endpoint(), nil, nil
	}
}
