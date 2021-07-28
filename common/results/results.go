package results

import (
	"runtime"

	"github.com/kataras/iris/v12"
)

// -------------------------------------------------------------------------
// 类型定义 Result: 内部返回值，例如 service, helper
// -------------------------------------------------------------------------

type Result struct {
	Success bool                   `json:"success"`
	Msg     string                 `json:"msg"`
	Data    map[string]interface{} `json:"data"`
}

// --------------------------------------------------------------------
// API：导出方法
// --------------------------------------------------------------------

func API() iris.Map {
	pc, _, _, _ := runtime.Caller(1)

	return iris.Map{
		"api_name":    runtime.FuncForPC(pc).Name() + "()",
		"api_expired": false,
	}
}

func API_expired() iris.Map {
	pc, _, _, _ := runtime.Caller(1)

	return iris.Map{
		"api_name":    runtime.FuncForPC(pc).Name() + "()",
		"api_expired": true,
	}
}

func Ok(data map[string]interface{}, msg ...string) *Result {

	var ok_data = data

	api, ok := data["api_name"]
	if ok == false {
		pc, _, _, _ := runtime.Caller(1)
		ok_data["api_name"] = runtime.FuncForPC(pc).Name() + "()"
	} else {
		ok_data["api_name"] = api
	}

	is_expired, ok := data["api_expired"]
	if ok == false {
		is_expired = false
		ok_data["api_expired"] = false
	} else {
		ok_data["api_expired"] = is_expired
	}

	r := &Result{Data: ok_data, Success: true}

	if len(msg) > 0 {
		r.Msg = msg[0]
	} else {
		r.Msg = "ok"
	}

	if is_expired == false {
		r.Success = true
	} else {
		r.Success = false
		r.Msg = "API is expired"
	}

	return r
}

func Error(data map[string]interface{}, msg ...string) *Result {

	var ok_data = data

	api, ok := data["api_name"]
	if ok == false {
		pc, _, _, _ := runtime.Caller(2)
		ok_data["api_name"] = runtime.FuncForPC(pc).Name() + "()"
	} else {
		ok_data["api_name"] = api
	}

	is_expired, ok := data["api_expired"]
	if ok == false {
		is_expired = false
		ok_data["api_expired"] = false
	} else {
		ok_data["api_expired"] = is_expired
	}

	r := &Result{Data: ok_data, Success: false}

	if len(msg) > 0 {
		r.Msg = msg[0]
	} else {
		r.Msg = "error"
	}

	if is_expired == true {
		r.Msg = "API is expired"
	}

	return r
}

func (r *Result) Is_success() bool {
	return r.Success == true
}

func (r *Result) Is_failure() bool {
	return r.Success == false
}

func (r *Result) Api_message() string {
	return r.Msg
}

func (r *Result) Api_name() string {
	return r.Data["api_name"].(string)
}

func (r *Result) Api_is_expired() bool {
	return r.Data["api_expired"].(bool) == true
}
