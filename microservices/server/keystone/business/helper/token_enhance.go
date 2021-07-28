package helper

import (
	"cctable/common/results"
	"cctable/microservices/server/keystone/business/model"
	"cctable/microservices/server/keystone/business/service"
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dgrijalva/jwt-go"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_token_enhance_helper interface {
	Enhance(ctx context.Context, token_uuid string) *results.Result
	Extrace(ctx context.Context, token_value string) *results.Result
}

type JWT_token_claims struct {
	User_name           string
	User_uuid           string
	Client_name         string
	Client_uuid         string
	Refresh_token_value string
	jwt.StandardClaims
}

type JWT_token_enhance struct {
	jwt_secret_key         []byte
	handler_token_store    service.I_token_store
	handler_user_service   service.I_user_service
	handler_client_service service.I_client_service
	handler_oauth_service  service.I_oauth_service
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_jwt_token_enhancer(db *mongo.Database, secret_key string) *JWT_token_enhance {
	token_enhancer := JWT_token_enhance{
		jwt_secret_key:         []byte(secret_key),
		handler_token_store:    service.New_token_store(db),
		handler_client_service: service.New_client_service(db),
		handler_oauth_service:  service.New_oauth_service(db),
		handler_user_service:   service.New_user_service(db),
	}

	return &token_enhancer
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (hlp *JWT_token_enhance) Enhance(ctx context.Context, token_uuid string) *results.Result {

	data := results.API()

	rs := hlp.handler_token_store.Get_token_by_uuid(ctx, token_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_token := rs.Data["token"].(model.Token)

	rs = hlp.handler_oauth_service.Get_oauth_by_uuid(ctx, the_token.Manager_oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_oauth := rs.Data["oauth"].(model.Oauth)

	rs = hlp.handler_user_service.Get_user_by_uuid(ctx, the_oauth.Manager_user_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_user := rs.Data["user"].(model.User)

	rs = hlp.handler_client_service.Get_client_by_uuid(ctx, the_oauth.Manager_client_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_client := rs.Data["client"].(model.Client)

	rs = hlp.handler_token_store.Get_refresh_token_by_oauth(ctx, the_oauth.Oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_refresh_token := rs.Data["token"].(model.Token)

	jwt_claims := JWT_token_claims{
		User_name:           the_user.User_name,
		User_uuid:           the_user.User_uuid,
		Client_name:         the_client.Client_name,
		Client_uuid:         the_client.Client_uuid,
		Refresh_token_value: the_refresh_token.Token_value,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: the_token.Expires_time.Unix(),
			Issuer:    "cctable-keystone",
		},
	}

	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt_claims)

	jwt_token_value, err := jwt_token.SignedString(hlp.jwt_secret_key)
	if err != nil {
		return results.Error(data, results.Err_failed_to_sign_jwt_token.Error())
	}

	rs = hlp.handler_token_store.Update_token_with_item(ctx, token_uuid, "token_type", model.TOKEN_type_jwt)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	rs = hlp.handler_token_store.Update_token_with_item(ctx, token_uuid, "token_value", jwt_token_value)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	data["token"] = rs.Data["token"]
	//pp.Println(data)

	return results.Ok(data)
}

func (hlp *JWT_token_enhance) Extrace(ctx context.Context, token_value string) *results.Result {

	data := results.API()

	jwt_token, err := jwt.ParseWithClaims(token_value, &JWT_token_claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return hlp.jwt_secret_key, nil
	})

	if err != nil {
		return results.Error(data, results.Err_failed_to_parse_jwt_token.Error())
	}

	claims := jwt_token.Claims.(*JWT_token_claims)

	rs := hlp.handler_oauth_service.Get_oauth_by_user_and_client_uuid(ctx, claims.User_uuid, claims.Client_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	the_oauth := rs.Data["oauth"].(model.Oauth)

	rs = hlp.handler_token_store.Get_access_token_by_oauth(ctx, the_oauth.Oauth_uuid)
	if rs.Is_failure() {
		return results.Error(data, rs.Msg)
	}

	data["token"] = rs.Data["token"]
	return results.Ok(data)
}
