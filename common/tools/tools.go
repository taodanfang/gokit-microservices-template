package tools

import (
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/kataras/iris/v12"
)

// --------------------------------------------------------------------
// 全局变量
// --------------------------------------------------------------------

var EmptyMap = iris.Map{}
var EmptyArrayString = make([]string, 0)
var EmptyArrayInterface = make([]interface{}, 0)

// --------------------------------------------------------------------
// API：导出方法
// --------------------------------------------------------------------

func NewUUID() string {
	a_uuid := uuid.Must(uuid.NewV4()).String()
	return a_uuid
}
