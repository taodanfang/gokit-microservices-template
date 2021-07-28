package main

import (
	"cctable/common/tools"
	"cctable/config"
	"cctable/gateway/reverse"
	"fmt"
	"github.com/rs/cors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.Init_config()
	cfg := config.Get_config()

	instance_manager := reverse.New_service_instance_reverse_manager()
	endpoint_manager := reverse.New_service_endpoint_reverse_manager(instance_manager)
	_ = endpoint_manager.Register_http_endpoint_reverse()

	err_ch := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		err_ch <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		tools.Log("Gateway reverse proxy start at: " + cfg.Gateway_proxy_reverse.Host + ":" + cfg.Gateway_proxy_reverse.Port)
		err_ch <- http.ListenAndServe(
			cfg.Gateway_proxy_reverse.Host+":"+cfg.Gateway_proxy_reverse.Port,
			cors.AllowAll().Handler(endpoint_manager.Reverse_proxy))
	}()

	tools.Log("exit", <-err_ch)
}
