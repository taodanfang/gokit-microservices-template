package transports

import (
	io "github.com/ambelovsky/gosf-socketio"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 2.1) 实现 application 接口

type Transport_for_ws_application struct {
	WS_server *io.Server
	DB        *mongo.Database
	Router    *iris.Application
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (t *Transport_for_ws_application) Make_route(route_name string, server *io.Server) {
	t.Router.Get(route_name, iris.FromStd(server))
}

func (t *Transport_for_ws_application) Get_db() *mongo.Database {
	return t.DB
}

func (t *Transport_for_ws_application) Get_router() *iris.Application {
	return t.Router
}

func (t *Transport_for_ws_application) Get_ws_server() *io.Server {
	return t.WS_server
}
