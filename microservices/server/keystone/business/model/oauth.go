package model

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// token 与 oauth 进行绑定
// oauth 记录 user/client 的授权情况

type Oauth struct {
	Oauth_uuid          string `json:"oauth_uuid", bson:"oauth_uuid"`
	Manager_user_uuid   string `json:"manager_user_uuid", bson:"manager_user_uuid"`
	Manager_client_uuid string `json:"manager_client_uuid", bson:"manager_client_uuid"`
}
