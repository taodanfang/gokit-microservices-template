package results

import "errors"

// --------------------------------------------------------------------
// 错误处理
// --------------------------------------------------------------------

// DB
var (
	Err_db_failed_to_find              = errors.New("错误：find 失败")
	Err_db_failed_to_insert            = errors.New("错误：insert 失败")
	Err_db_failed_to_update            = errors.New("错误：update 失败")
	Err_db_invalid_item_key_for_update = errors.New("错误：Update.item_key 错误")
)

// HTTP Request
var (
	Err_lost_authorization_with_request = errors.New("错误：无效请求（缺少Authorization）")
	Err_invalid_authorization_header    = errors.New("错误：无效请求（无效Authorization）")
	Err_invalid_client_request          = errors.New("错误：无效请求（未授权Oauth）")
	Err_invalid_token_request           = errors.New("错误：无效请求（令牌错AccessToken）")
	Err_lost_parameters_with_request    = errors.New("错误：无效请求（缺少参数）")
)

// GRPC Request
var (
	Err_connect_to_rpc_server_is_failure = errors.New("错误：连接GRPC服务失败")
	Err_invalid_method_with_rpc_request  = errors.New("错误：无效RPC调用(缺少method)")
	Err_invalid_params_with_rpc_request  = errors.New("错误：无效RPC调用(缺少params)")
)

// Service
var (
	Err_service_does_not_exist = errors.New("错误：服务不存在")
)

// User
var (
	Err_user_does_not_exist = errors.New("错误：用户不存在")
)

// Token
var (
	Err_token_does_not_exist            = errors.New("错误：Token 不存在")
	Err_token_has_expired               = errors.New("错误：Token 已经过期")
	Err_failed_to_sign_jwt_token        = errors.New("错误：JWT Token 签名失败")
	Err_failed_to_parse_jwt_token       = errors.New("错误：解析 JWT Token 失败")
	Err_does_not_support_the_grant_type = errors.New("错误：授权类型不支持")
	Err_invalid_user_name_and_password  = errors.New("错误：用户名和密码不正确")
	Err_invalid_refresh_token_value     = errors.New("错误：Refresh_token 不正确")
)

// Client
var (
	Err_client_does_not_exist = errors.New("错误：ClientID 不存在")
	Err_invalid_client_secret = errors.New("错误：ClientSecret 不正确")
)

// Oauth
var (
	Err_oauth_does_not_exist = errors.New("错误：Oauth 不存在")
)

// Room
var (
	Err_room_does_not_exist = errors.New("错误：用户名称不存在")
)

// Device
var (
	Err_device_does_not_exist = errors.New("错误：设备不存在")
)
