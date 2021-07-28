package main

import (
	"cctable/common/tools"
	"cctable/config"
	"cctable/gateway/proxy"
	"cctable/gateway/proxy/microserver"
	"github.com/kataras/iris/v12"
)

func main() {
	config.Init_config()
	cfg := config.Get_config()
	web_app := iris.New()

	instance_manager := proxy.New_service_instance_manager()
	endpoint_manager := proxy.New_service_endpoint_manager(instance_manager)

	microserver.Init_route_to_keystone(web_app, instance_manager, endpoint_manager)
	microserver.Init_route_to_message(web_app, instance_manager, endpoint_manager)

	tools.Log("Gateway server start at ", cfg.Gateway_proxy.Host+":"+cfg.Gateway_proxy.Port)
	tools.Log("err", web_app.Run(iris.Addr(cfg.Gateway_proxy.Host+":"+cfg.Gateway_proxy.Port)))
}
