package server_webs

import (
	"cctable/common/libs/gosf"
	"cctable/common/libs/micro/server/webs_server/transports"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

// 2.2) 创建 application

func Make_ws_application(router *iris.Application, db *mongo.Database) *transports.Transport_for_ws_application {

	tps := &transports.Transport_for_ws_application{Router: router, DB: db}

	transports.Init_plugin()

	gosf.Server_ready()
	tps.WS_server = gosf.Get_io_server()

	tps.Make_route("/connect", tps.WS_server)

	return tps
}
