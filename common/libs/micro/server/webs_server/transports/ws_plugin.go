package transports

import (
	"cctable/common/libs/gosf"
	"cctable/common/tools"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Plugin struct {
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func Init_plugin() {

	gosf.OnConnect(func(client *gosf.Client, request *gosf.Request) {

		tools.Log("connect a client.sid: ", client.Get_sid(), "client_uuid: ", client.Get_connecting_client_uuid())

		return
	})

	gosf.OnDisconnect(func(client *gosf.Client, request *gosf.Request) {

		tools.Log("disconnect a client.sid: ", client.Get_sid())

		return
	})

	gosf.OnBeforeRequest(func(client *gosf.Client, request *gosf.Request) {
	})

	gosf.OnAfterRequest(func(client *gosf.Client, request *gosf.Request, response *gosf.Message) {
	})

	gosf.OnBeforeResponse(func(client *gosf.Client, request *gosf.Request, response *gosf.Message) {
		//log.Println("Response for " + request.Endpoint + " endpoint is being prepared.")
	})
	gosf.OnAfterResponse(func(client *gosf.Client, request *gosf.Request, response *gosf.Message) {
		//log.Println("Response for " + request.Endpoint + " endpoint was sent.")
	})
	gosf.OnBeforeBroadcast(func(endpoint string, room string, response *gosf.Message) {
		//log.Println("Broadcast for " + endpoint + " endpoint is preparing to send to " + getRoom(room) + ".")
	})
	gosf.OnAfterBroadcast(func(endpoint string, room string, response *gosf.Message) {
		//log.Println("Broadcast for " + endpoint + " endpoint was sent to " + getRoom(room) + ".")
	})
	gosf.OnBeforeClientBroadcast(func(client *gosf.Client, endpoint string, room string, response *gosf.Message) {
		//log.Println("Broadcast for " + endpoint + " endpoint is preparing to send to " + getRoom(room) + ".")
	})
	gosf.OnAfterClientBroadcast(func(client *gosf.Client, endpoint string, room string, response *gosf.Message) {
		//log.Println("Broadcast for " + endpoint + " endpoint was sent to " + getRoom(room) + ".")
	})

	gosf.RegisterPlugin(new(Plugin))
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

// Activate is an aspect-oriented modular plugin requirement
func (p Plugin) Activate(app *gosf.AppSettings) {}

// Deactivate is an aspect-oriented modular plugin requirement
func (p Plugin) Deactivate(app *gosf.AppSettings) {}
