package model

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// client 代表将要访问服务的某个应用

type Client struct {
	Client_uuid                    string   `json:"client_uuid", bson:"client_uuid"`
	Client_name                    string   `json:"client_name", bson:"client_name"`
	Client_secret                  string   `json:"client_secret", bson:"client_secret"`
	Access_token_validity_seconds  int      `json:"access_token_validity_seconds", bson:"access_token_validity_seconds"`
	Refresh_token_validity_seconds int      `json:"refresh_token_validity_seconds", bson:"refresh_token_validity_seconds"`
	Authorized_grant_types         []string `json:"authorized_grant_types", bson:"authorized_grant_types"`
}
