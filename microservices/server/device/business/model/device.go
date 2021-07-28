package model

import "time"

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Device struct {
	Device_uuid string    `json:"device_uuid", bson:"device_uuid"`
	Device_name string    `json:"device_name", bson:"device_name"`
	Device_code string    `json:"device_code", bson:"device_code"`
	Create_at   time.Time `json:"create_at",bson:"create_at"`
}
