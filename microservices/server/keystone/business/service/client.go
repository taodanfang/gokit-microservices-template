package service

import (
	"cctable/common/results"
	"cctable/microservices/server/keystone/business/model"
	"context"
	"github.com/kataras/iris/v12"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

// 4.1) 定义单一服务接口

type I_client_service interface {
	New_client(ctx context.Context, client_name, client_secret string) *results.Result
	Update_client_with_item(ctx context.Context, client_uuid string, item_key string, item_value interface{}) *results.Result

	Get_client_by_uuid(ctx context.Context, client_uuid string) *results.Result
	Check_client_by_uuid_and_secret(ctx context.Context, client_uuid string, client_secret string) *results.Result
	Check_client_by_name_and_secret(ctx context.Context, client_name, client_secret string) *results.Result
}

// 4.2) 实现单一服务接口

type Client_service struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

// 4.3) 创建单一服务 service

func New_client_service(db *mongo.Database) *Client_service {
	service := Client_service{}
	service.db = db
	service.collection = db.Collection("clients")
	return &service
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *Client_service) New_client(ctx context.Context, client_name, client_secret string) *results.Result {

	data := results.API()

	a_client := model.Client{
		Client_uuid:   primitive.NewObjectID().Hex(),
		Client_name:   client_name,
		Client_secret: client_secret,
	}

	_, err := s.collection.InsertOne(ctx, a_client)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_insert.Error())
	}

	current := s.get_client(ctx, a_client.Client_uuid)
	data["client"] = current["client"]

	return results.Ok(data)
}

func (s *Client_service) Get_client_by_uuid(ctx context.Context, client_uuid string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"client_uuid", client_uuid},
	}

	var a_client = model.Client{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_client)

	if err != nil {
		return results.Error(data, results.Err_client_does_not_exist.Error())
	}

	data["client"] = a_client

	return results.Ok(data)
}

func (s *Client_service) Check_client_by_uuid_and_secret(ctx context.Context, client_uuid string, client_secret string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"client_uuid", client_uuid},
	}

	var a_client = model.Client{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_client)

	if err != nil {
		//pp.Println(err)
		return results.Error(data, results.Err_client_does_not_exist.Error())
	}

	if a_client.Client_secret != client_secret {
		return results.Error(data, results.Err_invalid_client_secret.Error())
	}

	data["client"] = a_client
	return results.Ok(data)
}

func (s *Client_service) Check_client_by_name_and_secret(ctx context.Context, client_name, client_secret string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"client_name", client_name},
		{"client_secret", client_secret},
	}

	var a_client = model.Client{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_client)

	if err != nil {
		return results.Error(data, results.Err_client_does_not_exist.Error())
	}

	data["client"] = a_client

	return results.Ok(data)
}

func (s *Client_service) Update_client_with_item(ctx context.Context, client_uuid string, item_key string, item_value interface{}) *results.Result {

	data := results.API()

	filter := bson.D{
		{"client_uuid", client_uuid},
	}

	update := bson.D{{}}

	switch item_key {
	case "client_name":
		update = bson.D{{"$set", bson.D{
			{item_key, item_value.(string)},
		}}}
	case "client_secret":
		update = bson.D{{"$set", bson.D{
			{item_key, item_value.(string)},
		}}}
	case "access_token_validity_seconds":
		update = bson.D{{"$set", bson.D{
			{item_key, item_value.(int)},
		}}}
	case "refresh_token_validity_seconds":
		update = bson.D{{"$set", bson.D{
			{item_key, item_value.(int)},
		}}}
	case "authorized_grant_types":
		update = bson.D{{"$set", bson.D{
			{item_key, item_value.([]string)},
		}}}

	default:
		return results.Error(data, results.Err_db_invalid_item_key_for_update.Error())
	}

	var new_client = model.Client{}
	err := s.collection.FindOneAndUpdate(ctx, filter, update).Decode(&new_client)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_update.Error())
	}

	current := s.get_client(ctx, client_uuid)
	data["client"] = current["client"]

	return results.Ok(data)
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (s *Client_service) get_client(ctx context.Context, client_uuid string) (data iris.Map) {

	filter := bson.D{{"client_uuid", client_uuid}}

	var a_client = model.Client{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_client)

	data = iris.Map{
		"client": iris.Map{},
	}

	if err == nil {
		data = iris.Map{
			"client": a_client,
		}
	}

	return data
}
