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

type I_oauth_service interface {
	New_oauth(ctx context.Context, user_uuid, client_uuid string) *results.Result
	Get_oauth_by_uuid(ctx context.Context, oauth_uuid string) *results.Result
	Get_oauth_by_user_and_client_uuid(ctx context.Context, user_uuid, client_uuid string) *results.Result
}

type Oauth_service struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_oauth_service(db *mongo.Database) *Oauth_service {
	service := Oauth_service{}
	service.db = db
	service.collection = db.Collection("oauths")
	return &service
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *Oauth_service) New_oauth(ctx context.Context, user_uuid, client_uuid string) *results.Result {

	data := results.API()

	a_oauth := model.Oauth{
		Oauth_uuid:          primitive.NewObjectID().Hex(),
		Manager_user_uuid:   user_uuid,
		Manager_client_uuid: client_uuid,
	}

	_, err := s.collection.InsertOne(ctx, a_oauth)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_insert.Error())
	}

	current := s.get_oauth(ctx, a_oauth.Oauth_uuid)
	data["oauth"] = current["oauth"]

	return results.Ok(data)
}

func (s *Oauth_service) Get_oauth_by_uuid(ctx context.Context, oauth_uuid string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"oauth_uuid", oauth_uuid},
	}

	//pp.Println(filter)

	var a_oauth = model.Oauth{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_oauth)

	if err != nil {
		return results.Error(data, results.Err_oauth_does_not_exist.Error())
	}

	data["oauth"] = a_oauth

	return results.Ok(data)
}

func (s *Oauth_service) Get_oauth_by_user_and_client_uuid(ctx context.Context, user_uuid string, client_uuid string) *results.Result {

	data := results.API()

	filter := bson.D{
		{"manager_user_uuid", user_uuid},
		{"manager_client_uuid", client_uuid},
	}

	//pp.Println(filter)

	var a_oauth = model.Oauth{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_oauth)

	if err != nil {
		//pp.Println(err)
		return results.Error(data, results.Err_oauth_does_not_exist.Error())
	}

	data["oauth"] = a_oauth

	return results.Ok(data)
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (s *Oauth_service) get_oauth(ctx context.Context, oauth_uuid string) (data iris.Map) {

	filter := bson.D{{"oauth_uuid", oauth_uuid}}

	var a_oauth = model.Oauth{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_oauth)

	data = iris.Map{
		"oauth": iris.Map{},
	}

	if err == nil {
		data = iris.Map{
			"oauth": a_oauth,
		}
	}

	return data
}
