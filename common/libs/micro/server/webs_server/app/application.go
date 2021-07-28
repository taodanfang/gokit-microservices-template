package app

import (
	io "github.com/ambelovsky/gosf-socketio"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 1.1) 定义 application 接口

type I_ws_application interface {
	Make_route(string, *io.Server)

	Get_db() *mongo.Database
	Get_router() *iris.Application

	Get_ws_server() *io.Server
}
