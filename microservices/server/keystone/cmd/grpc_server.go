package main

import (
	"cctable/common/db"
	"cctable/common/libs/micro/discovery/discover_grpc"
	"cctable/common/libs/micro/server/grpc_server/transports"
	"cctable/common/tools"
	"cctable/microservices/server/keystone/server_grpc/pb_server"
	"cctable/microservices/server/keystone/server_grpc/proto"
	"fmt"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os/signal"
	"syscall"

	"cctable/config"
	"cctable/microservices/server/keystone/server_grpc/endpoints"
	go_kit_log "github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"os"
)

func main() {
	config.Init_config()
	cfg := config.Get_config()
	discovery_client, err := discover_grpc.New_kit_discovery_grpc_client(cfg.Consul.Host, cfg.Consul.Port)
	if err != nil {
		log.Println("Get consul client failed!")
		os.Exit(-1)
	}

	err_chan := make(chan error)
	db.Init_mongo_db()
	database := db.Get_mongo_db()

	logger := go_kit_log.NewLogfmtLogger(os.Stderr)
	logger = go_kit_log.With(logger, "ts", go_kit_log.DefaultTimestampUTC)
	logger = go_kit_log.With(logger, "caller", go_kit_log.DefaultCaller)

	grpc_server := transports.Make_grpc_application(database, logger)
	endpoints.Init_grpc_endpoint_managers(grpc_server)
	router := pb_server.New_token_service_grpc_router(grpc_server)

	instance_id := cfg.Keystone.Name + "-grpc-" + cfg.Keystone.ID

	if !discovery_client.Register(cfg.Keystone.Name, instance_id, cfg.Keystone.Host, cfg.Keystone.Grpc_port) {
		tools.Log("register grpc service: " + cfg.Keystone.Name + " failed.")
		os.Exit(-1)
	}

	go func() {
		tools.Log(cfg.Keystone.Name + " grpc server start at " + cfg.Keystone.Host + ":" + cfg.Keystone.Grpc_port)

		ls, _ := net.Listen("tcp", cfg.Keystone.Host+":"+cfg.Keystone.Grpc_port)
		gRPCServer := grpc.NewServer()
		proto.RegisterTokenServiceServer(gRPCServer, router)
		grpc_health_v1.RegisterHealthServer(gRPCServer, &discover_grpc.Gprc_discovery_service{})
		err_chan <- gRPCServer.Serve(ls)
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
