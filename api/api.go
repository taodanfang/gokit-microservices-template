package api

// --------------------------------------------------------------------
// 服务列表
// --------------------------------------------------------------------

const (
	MS_http_server_keystone = "keystone"
	MS_http_server_message  = "message"

	MS_grpc_server_keystone = "keystone"
)

// --------------------------------------------------------------------
// HTTP API
// --------------------------------------------------------------------

const (
	EDP_http_keystone__oauth__grant_token  = "/oauth/grant_token"
	EDP_http_keystone__oauth__check_token  = "/oauth/check_token"
	EDP_http_keystone__user__register_user = "/user/register_user"
	EDP_http_keystone__user__login         = "/user/login"
	EDP_http_keystone__user__get_all_users = "/user/get_all_users"
	EDP_http_keystone__user__logout        = "/user/logout"
)

const (
	EDP_http_message__room__check_room = "/room/check_room"
)

const (
	EDP_http_device__device__register_device = "/device/register_device"
	EDP_http_device__device__get_all_devices = "/device/get_all_devices"
)

// --------------------------------------------------------------------
// GRPC API ( 模仿http，但不以'/'开始 ）
// --------------------------------------------------------------------

const (
	SVC_grpc_keystone_pb_token_service = "TokenService"
)

const (
	EDP_grpc_keystone_token_check_access_token = "/token/check_access_token"
)
