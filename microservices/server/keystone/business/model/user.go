package model

import "time"

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type User struct {
	User_uuid string    `json:"user_uuid", bson:"user_uuid"`
	User_name string    `json:"user_name", bson:"user_name"`
	Telephone string    `json:"telephone", bson:"telephone"`
	Password  string    `json:"password", bson:"password"`
	Actors    []string  `json:"actors", bson:"actors"`
	Create_at time.Time `json:"create_at",bson:"create_at"`

	Is_login bool `json:"is_login",bson:"is_login"`
}
