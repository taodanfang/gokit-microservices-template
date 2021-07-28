package main

import (
	"cctable/common/db"
	"cctable/common/libs/micro/discovery/discover_http"
	"cctable/common/libs/micro/server/http_server/transports"
	"cctable/common/tools"
	"cctable/config"
	"cctable/microservices/server/device/server_http/endpoints"

	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/iris/v12"
)

func main() {

	config.Init_config()
	cfg := config.Get_config()
	discovery_client, err := discover_http.New_kit_discovery_http_client(cfg.Consul.Host, cfg.Consul.Port)
	if err != nil {
		log.Println("Get consul client failed!")
		os.Exit(-1)
	}

	err_chan := make(chan error)
	db.Init_mongo_db()

	database := db.Get_mongo_db()
	web_app := iris.New()

	discover_http.Make_http_discovery(web_app, discovery_client)
	http_server := transports.Make_http_application(web_app, database)
	endpoints.Init_http_endpoint_managers(http_server)

	instance_id := cfg.Device.Name + "-http-" + cfg.Device.ID

	if !discovery_client.Register(cfg.Device.Name, instance_id, cfg.Device.Host, cfg.Device.Http_port) {
		tools.Log("register http service: " + cfg.Device.Name + " failed.")
		os.Exit(-1)
	}

	go func() {
		tools.Log(cfg.Device.Name + " http server start at " + cfg.Device.Host + ":" + cfg.Device.Http_port)
		err_chan <- web_app.Run(iris.Addr(cfg.Device.Host + ":" + cfg.Device.Http_port))
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		err_chan <- fmt.Errorf("%s", <-c)
	}()

	err = <-err_chan

	discovery_client.Deregister(instance_id)
	tools.Log(err)
}
