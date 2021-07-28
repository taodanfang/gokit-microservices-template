package server_http

import (
	"cctable/common/libs/micro/server/http_server/transports"
	"cctable/common/tools"
	"github.com/go-kit/kit/transport"

	"go.mongodb.org/mongo-driver/mongo"

	go_kit_http "github.com/go-kit/kit/transport/http"
	"github.com/kataras/iris/v12"
)

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

// 2.2) 创建 application

func Make_http_application(router *iris.Application, db *mongo.Database) *transports.Transport_for_http_application {

	tps := &transports.Transport_for_http_application{Router: router, DB: db}

	client_authorization_options := []go_kit_http.ServerOption{
		go_kit_http.ServerErrorHandler(transport.NewLogErrorHandler(tools.Get_gokit_logger())),
		go_kit_http.ServerErrorEncoder(tps.Encode_error),
	}

	tps.Server_options = client_authorization_options

	return tps
}
