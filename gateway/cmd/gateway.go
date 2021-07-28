package main

import (
	"cctable/common/tools"
	"cctable/config"
	"cctable/gateway/proxy"
	"cctable/gateway/proxy/microserver"
	"cctable/gateway/reverse"
	"fmt"
	"github.com/kataras/iris/v12"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.Init_config()
	cfg := config.Get_config()
	web_app := iris.New()

	reverse_instance_manager := reverse.New_service_instance_reverse_manager()
	reverse_endpoint_manager := reverse.New_service_endpoint_reverse_manager(reverse_instance_manager)
	_ = reverse_endpoint_manager.Register_http_endpoint_reverse()

	proxy_instance_manager := proxy.New_service_instance_manager()
	proxy_endpoint_manager := proxy.New_service_endpoint_manager(proxy_instance_manager)
	microserver.Init_route_to_keystone(web_app, proxy_instance_manager, proxy_endpoint_manager)
	microserver.Init_route_to_message(web_app, proxy_instance_manager, proxy_endpoint_manager)

	err_ch := make(chan error)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		err_ch <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		tools.Log("Gateway http server start at " + cfg.Gateway_proxy.Host + ":" + cfg.Gateway_proxy.Port)
		err_ch <- web_app.Run(iris.Addr(cfg.Gateway_proxy.Host + ":" + cfg.Gateway_proxy.Port))
	}()

	go func() {
		tools.Log("Gateway websocket server start at " + cfg.Gateway_proxy_reverse.Host + ":" + cfg.Gateway_proxy_reverse.Port)
		err_ch <- http.ListenAndServe(cfg.Gateway_proxy_reverse.Host+":"+cfg.Gateway_proxy_reverse.Port, reverse_endpoint_manager.Reverse_proxy)
	}()

	tools.Log("exit", <-err_ch)
}
