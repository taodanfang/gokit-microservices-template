package main

import (
	"cctable/common/db"
	"cctable/common/libs/micro/discovery/discover_http"
	"cctable/common/tools"
	"cctable/config"
	"cctable/microservices/server/message/business/service"
	"cctable/microservices/server/message/server_http"
	http_endpoints "cctable/microservices/server/message/server_http/endpoints"
	"cctable/microservices/server/message/server_webs"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"

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

	ctx := context.Background()
	database := db.Get_mongo_db()
	web_app := iris.New()

	init_data(ctx, database)
	discover_http.Make_http_discovery(web_app, discovery_client)

	http_server := server_http.Make_http_application(web_app, database)
	http_endpoints.Init_http_endpoint_managers(http_server)

	server_webs.Make_ws_application(web_app, database)

	instance_id := cfg.Message.Name + "-" + cfg.Message.ID

	go func() {
		tools.Log("Http server start at " + cfg.Message.Host + ":" + cfg.Message.Http_port)

		if !discovery_client.Register(
			cfg.Message.Name, instance_id, "/health",
			cfg.Message.Host, cfg.Message.Http_port, nil) {
			tools.Log("service: " + cfg.Message.Name + "failed.")
			os.Exit(-1)
		}

		err_chan <- web_app.Run(iris.Addr(cfg.Message.Host + ":" + cfg.Message.Http_port))
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		err_chan <- fmt.Errorf("%s", <-c)
	}()

	err = <-err_chan
	log.Println(err)
}

func init_data(ctx context.Context, database *mongo.Database) {

	var room_service service.I_room_service

	room_service = service.New_room_service(database)

	rs := room_service.Check_room_by_name(ctx, "test-room-1")
	if rs.Is_failure() {
		rs = room_service.New_room(ctx, "test-room-1")
	}
}
