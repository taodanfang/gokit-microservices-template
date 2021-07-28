package main

import (
	"cctable/common/db"
	"cctable/common/libs/micro/discovery/discover_http"
	"cctable/common/libs/micro/server/http_server/transports"
	"cctable/common/tools"
	"cctable/config"
	"cctable/microservices/server/keystone/business/model"
	"cctable/microservices/server/keystone/business/service"
	"cctable/microservices/server/keystone/server_http/endpoints"
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
	http_server := transports.Make_http_application(web_app, database)
	endpoints.Init_http_endpoint_managers(http_server)

	instance_id := cfg.Keystone.Name + "-http-" + cfg.Keystone.ID

	if !discovery_client.Register(cfg.Keystone.Name, instance_id, cfg.Keystone.Host, cfg.Keystone.Http_port) {
		tools.Log("register http service: " + cfg.Keystone.Name + " failed.")
		os.Exit(-1)
	}

	go func() {
		tools.Log(cfg.Keystone.Name + " http server start at " + cfg.Keystone.Host + ":" + cfg.Keystone.Http_port)
		err_chan <- web_app.Run(iris.Addr(cfg.Keystone.Host + ":" + cfg.Keystone.Http_port))
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

func init_data(ctx context.Context, database *mongo.Database) {

	var client_service service.I_client_service
	var user_service service.I_user_service
	var oauth_service service.I_oauth_service

	client_service = service.New_client_service(database)
	user_service = service.New_user_service(database)
	oauth_service = service.New_oauth_service(database)

	var client_uuid string
	rs := client_service.Check_client_by_name_and_secret(ctx, "cctable-web-client", "cctable-web-client-secret")
	if rs.Is_failure() {
		rs = client_service.New_client(ctx, "cctable-web-client", "cctable-web-client-secret")
	}
	client_uuid = rs.Data["client"].(model.Client).Client_uuid

	_ = client_service.Update_client_with_item(ctx, client_uuid, "access_token_validity_seconds", 3600)
	_ = client_service.Update_client_with_item(ctx, client_uuid, "refresh_token_validity_seconds", 3600*12)
	_ = client_service.Update_client_with_item(ctx, client_uuid, "authorized_grant_types", []string{"password", "refresh_token"})

	var simple_user_uuid string
	rs = user_service.Check_user_by_name_and_password(ctx, "simple", "123456")
	if rs.Is_failure() {
		rs = user_service.New_user(ctx, "simple", "123456")
	}
	simple_user_uuid = rs.Data["user"].(model.User).User_uuid

	rs = oauth_service.Get_oauth_by_user_and_client_uuid(ctx, simple_user_uuid, client_uuid)
	if rs.Is_failure() {
		_ = oauth_service.New_oauth(ctx, simple_user_uuid, client_uuid)
	}

	var admin_user_uuid string
	rs = user_service.Check_user_by_name_and_password(ctx, "admin", "123456")
	if rs.Is_failure() {
		rs = user_service.New_user(ctx, "admin", "123456")
	}
	admin_user_uuid = rs.Data["user"].(model.User).User_uuid

	rs = oauth_service.Get_oauth_by_user_and_client_uuid(ctx, admin_user_uuid, client_uuid)
	if rs.Is_failure() {
		_ = oauth_service.New_oauth(ctx, admin_user_uuid, client_uuid)
	}
}
