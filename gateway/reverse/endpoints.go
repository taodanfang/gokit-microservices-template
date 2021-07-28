package reverse

import (
	"cctable/common/tools"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Service_endpoint_reverse_manager struct {
	instance_manager *Service_instance_reverse_manager

	Reverse_proxy *httputil.ReverseProxy
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_service_endpoint_reverse_manager(instance_manager *Service_instance_reverse_manager) *Service_endpoint_reverse_manager {

	return &Service_endpoint_reverse_manager{
		instance_manager: instance_manager,
	}
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *Service_endpoint_reverse_manager) Register_http_endpoint_reverse() error {
	director := func(req *http.Request) {
		req_path := req.URL.Path
		if req_path == "" {
			return
		}

		//pp.Println("req_pqth: ", req_path)

		path_array := strings.Split(req_path, "/")
		service_name := path_array[1]

		//pp.Println("service_name: ", service_name)
		//tools.Log(service_name)

		result, _, err := s.instance_manager.api_client.Catalog().Service(service_name, "http", nil)
		if err != nil {
			tools.Log(err)
			return
		}

		if len(result) == 0 {
			tools.Log("no such service instance: ", service_name)
			return
		}

		dest_path := strings.Join(path_array[2:], "/")

		//tools.Log("result: ", result)

		tgt := result[rand.Int()%len(result)]

		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d", tgt.ServiceAddress, tgt.ServicePort)
		req.URL.Path = "/" + dest_path

		//tools.Log(req.URL.Host, req.URL.Path)

	}

	the_proxy := &httputil.ReverseProxy{Director: director}

	s.Reverse_proxy = the_proxy

	return nil
}
