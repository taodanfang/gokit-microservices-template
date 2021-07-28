package model

import (
	"time"
)

// --------------------------------------------------------------------
// 常量定义
// --------------------------------------------------------------------

const (
	TOKEN_type_uuid = "UUID"
	TOKEN_type_jwt  = "JWT"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Token struct {
	Token_uuid   string    `json:"token_uuid", bson:"token_uuid"`
	Token_type   string    `json:"token_type", bson:"token_type"`
	Token_value  string    `json:"token_value", bson:"token_value"`
	Expires_time time.Time `json:"expires_time", bson:"expires_time"`

	Manager_oauth_uuid string `json:"manager_oauth_uuid", bson:"manager_oauth_uuid"`
	Is_refresh_token   bool   `json:"is_refresh_token", bson:"is_refresh_token"`
}

// --------------------------------------------------------------------
// 导出方法
// --------------------------------------------------------------------

func (t *Token) Is_expired() bool {
	return t.Expires_time.IsZero() != true && t.Expires_time.Before(time.Now())
}
