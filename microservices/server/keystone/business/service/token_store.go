package service

import (
	"cctable/common/results"
	"cctable/common/tools"
	"cctable/microservices/server/keystone/business/model"
	"context"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_token_store interface {
	Store_access_token(ctx context.Context, token_type, token_value string, expires_time time.Time, manager_oauth_uuid string) *results.Result
	Read_access_token(ctx context.Context, token_value string) *results.Result
	Get_access_token_by_oauth(ctx context.Context, manager_oauth_uuid string) *results.Result
	Remove_access_token(ctx context.Context, token_value string) *results.Result

	Store_refresh_token(ctx context.Context, token_type, token_value string, expires_time time.Time, manager_oauth_uuid string) *results.Result
	Read_refresh_token(ctx context.Context, token_value string) *results.Result
	Get_refresh_token_by_oauth(ctx context.Context, manager_oauth_uuid string) *results.Result
	Remove_refresh_token(ctx context.Context, token_value string) *results.Result

	Update_token_with_item(ctx context.Context, token_uuid string, item_key string, item_value interface{}) *results.Result
	Get_token_by_uuid(ctx context.Context, token_uuid string) *results.Result
}

type Token_store struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_token_store(db *mongo.Database) *Token_store {
	store := Token_store{}
	store.db = db
	store.collection = db.Collection("tokens")
	return &store
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *Token_store) Get_token_by_uuid(ctx context.Context, token_uuid string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"token_uuid", token_uuid},
	}

	var a_token = model.Token{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_token)

	if err != nil {
		return results.Error(data, results.Err_token_does_not_exist.Error())
	}

	data["token"] = a_token

	return results.Ok(data)
}

func (s *Token_store) Read_access_token(ctx context.Context, token_value string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"token_value", token_value},
		{"is_refresh_token", false},
	}

	//tools.Log(filter)

	var a_token = model.Token{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_token)

	if err != nil {
		//tools.Log(err)
		return results.Error(data, results.Err_token_does_not_exist.Error())
	}

	data["token"] = a_token

	return results.Ok(data)
}

func (s *Token_store) Update_token_with_item(ctx context.Context, token_uuid string, item_key string, item_value interface{}) *results.Result {

	data := results.API()

	filter := bson.D{
		{"token_uuid", token_uuid},
	}

	update := bson.D{{}}

	switch item_key {
	case "token_type":
		update = bson.D{{"$set", bson.D{
			{item_key, item_value.(string)},
		}}}
	case "token_value":
		update = bson.D{{"$set", bson.D{
			{item_key, item_value.(string)},
		}}}
	case "expires_time":
		update = bson.D{{"$set", bson.D{
			{item_key, item_value.(time.Time)},
		}}}
	case "manager_oauth_uuid":
		update = bson.D{{"$set", bson.D{
			{item_key, item_value.(string)},
		}}}
	case "is_refresh_token":
		update = bson.D{{"$set", bson.D{
			{item_key, item_value.(bool)},
		}}}

	default:
		return results.Error(data, results.Err_db_invalid_item_key_for_update.Error())
	}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_update.Error())
	}

	current := s.get_token(ctx, token_uuid)
	data["token"] = current["token"]

	return results.Ok(data)
}

func (s *Token_store) Store_access_token(ctx context.Context, token_type, token_value string, expires_time time.Time, manager_oauth_uuid string) *results.Result {

	data := results.API()

	a_token := model.Token{
		Token_uuid:         primitive.NewObjectID().Hex(),
		Token_type:         token_type,
		Token_value:        token_value,
		Expires_time:       expires_time,
		Manager_oauth_uuid: manager_oauth_uuid,
		Is_refresh_token:   false,
	}

	_, err := s.collection.InsertOne(ctx, a_token)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_insert.Error())
	}

	current := s.get_token(ctx, a_token.Token_uuid)
	data["token"] = current["token"]

	return results.Ok(data)
}

func (s *Token_store) Get_access_token_by_oauth(ctx context.Context, manager_oauth_uuid string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"is_refresh_token", false},
		{"manager_oauth_uuid", manager_oauth_uuid},
	}

	var a_token = model.Token{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_token)

	if err != nil {
		return results.Error(data, results.Err_token_does_not_exist.Error())
	}

	data["token"] = a_token

	return results.Ok(data)
}

func (s *Token_store) Remove_access_token(ctx context.Context, token_value string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"token_value", token_value},
		{"is_refresh_token", false},
	}

	rst := s.collection.FindOneAndDelete(ctx, filter)

	if rst.Err() != nil {
		tools.Log(rst.Err())
		return results.Error(data, results.Err_token_does_not_exist.Error())
	}

	return results.Ok(data)
}

func (s *Token_store) Store_refresh_token(ctx context.Context, token_type, token_value string, expires_time time.Time, manager_oauth_uuid string) *results.Result {

	data := results.API()

	a_token := model.Token{
		Token_uuid:         primitive.NewObjectID().Hex(),
		Token_type:         token_type,
		Token_value:        token_value,
		Expires_time:       expires_time,
		Manager_oauth_uuid: manager_oauth_uuid,
		Is_refresh_token:   true,
	}

	_, err := s.collection.InsertOne(ctx, a_token)
	if err != nil {

		return results.Error(data, results.Err_db_failed_to_insert.Error())
	}

	current := s.get_token(ctx, a_token.Token_uuid)
	data["token"] = current["token"]

	return results.Ok(data)
}

func (s *Token_store) Read_refresh_token(ctx context.Context, token_value string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"token_value", token_value},
		{"is_refresh_token", true},
	}

	var a_token = model.Token{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_token)

	if err != nil {
		return results.Error(data, results.Err_token_does_not_exist.Error())
	}

	data["token"] = a_token

	return results.Ok(data)
}

func (s *Token_store) Get_refresh_token_by_oauth(ctx context.Context, manager_oauth_uuid string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"is_refresh_token", true},
		{"manager_oauth_uuid", manager_oauth_uuid},
	}

	var a_token = model.Token{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_token)

	if err != nil {
		return results.Error(data, results.Err_token_does_not_exist.Error())
	}

	data["token"] = a_token

	return results.Ok(data)
}

func (s *Token_store) Remove_refresh_token(ctx context.Context, token_value string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"token_value", token_value},
		{"is_refresh_token", true},
	}

	rst := s.collection.FindOneAndDelete(ctx, filter)

	if rst.Err() != nil {
		tools.Log(rst.Err())
		return results.Error(data, results.Err_token_does_not_exist.Error())
	}

	return results.Ok(data)
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (s *Token_store) get_token(ctx context.Context, token_uuid string) (data iris.Map) {

	filter := bson.D{{"token_uuid", token_uuid}}

	var a_token = model.Token{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_token)

	data = iris.Map{
		"token": iris.Map{},
	}

	if err == nil {
		data = iris.Map{
			"token": a_token,
		}
	}

	return data
}
