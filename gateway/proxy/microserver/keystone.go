package microserver

import (
	"cctable/api"
	"cctable/gateway/proxy"
	"github.com/kataras/iris/v12"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Proxy_to_keystone struct {
	router           *iris.Application
	service_name     string
	transport_type   string
	instance_manager *proxy.Service_instance_manager
	endpoint_manager *proxy.Service_endpoint_manager
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func Init_route_to_keystone(router *iris.Application, instance_manager *proxy.Service_instance_manager, endpoint_manager *proxy.Service_endpoint_manager) {
	an_proxy := Proxy_to_keystone{
		service_name:     api.MS_http_server_keystone,
		transport_type:   "http",
		router:           router,
		instance_manager: instance_manager,
		endpoint_manager: endpoint_manager}

	instance_manager.Register_instance(api.MS_http_server_keystone)

	_ = an_proxy.Make_route("POST", api.EDP_http_keystone__oauth__grant_token)
}

func (p *Proxy_to_keystone) Make_route(method, path string) error {
	edp := p.endpoint_manager
	err := edp.Register_http_endpoint(p.service_name, path, method)
	if err != nil {
		return err
	}

	the_handler, _ := edp.Get_handler(p.service_name, path)

	//pp.Println("handler: ", the_handler)

	full_path := p.service_name + "/" + path
	switch method {
	case "GET":
		p.router.Get(full_path, iris.FromStd(the_handler))
	case "POST":
		p.router.Post(full_path, iris.FromStd(the_handler))
	}

	//pp.Println("full_path: ", full_path)
	return nil
}
