package service

import (
	"cctable/common/results"
	"cctable/common/tools"
	"cctable/microservices/server/keystone/business/model"
	"context"
	"time"

	"github.com/kataras/iris/v12"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_user_service interface {
	New_user(ctx context.Context, user_name, password string) *results.Result
	New_user_with_telephone(ctx context.Context, user_name, telephone string) *results.Result

	Update_user(ctx context.Context, user_uuid string, item_key string, item_value interface{}) *results.Result

	Check_user_by_name_and_telephone(ctx context.Context, user_name, telephone string) *results.Result
	Check_user_by_name_and_password(ctx context.Context, user_name, password string) *results.Result

	Get_user_by_uuid(ctx context.Context, user_uuid string) *results.Result

	Get_all_users(ctx context.Context) *results.Result
}

type User_service struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_user_service(db *mongo.Database) *User_service {
	service := User_service{}
	service.db = db
	service.collection = db.Collection("users")
	return &service
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *User_service) New_user(ctx context.Context, user_name, password string) *results.Result {
	data := results.API()

	a_user := model.User{
		User_uuid: primitive.NewObjectID().Hex(),
		User_name: user_name,
		Password:  password,
		Create_at: time.Now(),
		Actors:    tools.EmptyArrayString,
		Is_login:  false,
	}

	_, err := s.collection.InsertOne(ctx, a_user)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_insert.Error())
	}

	current := s.get_user(ctx, a_user.User_uuid)
	data["user"] = current["user"]

	return results.Ok(data)
}

func (s *User_service) New_user_with_telephone(ctx context.Context, user_name, telephone string) *results.Result {
	data := results.API()

	a_user := model.User{
		User_uuid: primitive.NewObjectID().Hex(),
		User_name: user_name,
		Password:  telephone,
		Telephone: telephone,
		Create_at: time.Now(),
		Actors:    tools.EmptyArrayString,
		Is_login:  false,
	}

	_, err := s.collection.InsertOne(ctx, a_user)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_insert.Error())
	}

	current := s.get_user(ctx, a_user.User_uuid)
	data["user"] = current["user"]

	return results.Ok(data)
}

func (s *User_service) Update_user(ctx context.Context, user_uuid string, item_key string, item_value interface{}) *results.Result {
	data := results.API()

	filter := bson.D{
		{"user_uuid", user_uuid},
	}

	update := bson.D{{}}

	switch item_key {
	case "user_name":
		update = bson.D{{"$set", bson.D{{item_key, item_value.(string)}}}}
		break
	case "telephone":
		update = bson.D{{"$set", bson.D{{item_key, item_value.(string)}}}}
		break
	case "password":
		update = bson.D{{"$set", bson.D{{item_key, item_value.(string)}}}}
		break
	case "is_login":
		update = bson.D{{"$set", bson.D{{item_key, item_value.(bool)}}}}
		break
	default:
		return results.Error(data, results.Err_db_invalid_item_key_for_update.Error())
	}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		tools.Log(err)
		return results.Error(data, results.Err_db_failed_to_update.Error())
	}

	current := s.get_user(ctx, user_uuid)
	data["user"] = current["user"]

	return results.Ok(data)
}

func (s *User_service) Get_user_by_uuid(ctx context.Context, user_uuid string) *results.Result {
	data := results.API()

	filter := bson.D{
		{"user_uuid", user_uuid},
	}

	var a_user = model.User{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_user)

	if err != nil {
		return results.Error(data, results.Err_user_does_not_exist.Error())
	}

	data["user"] = a_user

	return results.Ok(data)
}

func (s *User_service) Check_user_by_name_and_telephone(ctx context.Context, user_name, telephone string) *results.Result {
	data := results.API()

	filter := bson.D{
		{"user_name", user_name},
		{"telephone", telephone},
	}

	var a_user = model.User{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_user)

	if err != nil {
		return results.Error(data, results.Err_user_does_not_exist.Error())
	}

	data["user"] = a_user

	return results.Ok(data)
}

func (s *User_service) Check_user_by_name_and_password(ctx context.Context, user_name, password string) *results.Result {
	data := results.API()

	filter := bson.D{
		{"user_name", user_name},
		{"password", password},
	}

	var a_user = model.User{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_user)

	if err != nil {
		return results.Error(data, results.Err_user_does_not_exist.Error())
	}

	data["user"] = a_user

	return results.Ok(data)
}

func (s *User_service) Get_all_users(ctx context.Context) *results.Result {
	data := results.API()

	filter := bson.D{{}}

	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_find.Error())
	}

	var users []interface{}
	for cur.Next(ctx) {
		var a_user model.User
		err := cur.Decode(&a_user)
		if err != nil {
			return results.Error(data, results.Err_db_failed_to_find.Error())
		}
		users = append(users, a_user)
	}

	_ = cur.Close(ctx)

	if len(users) <= 0 {
		data["users"] = tools.EmptyArrayString
	} else {
		data["users"] = users
	}

	return results.Ok(data)
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (s *User_service) get_user(ctx context.Context, user_uuid string) (data iris.Map) {

	filter := bson.D{{"user_uuid", user_uuid}}

	var a_user = model.User{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_user)

	data = iris.Map{
		"user": iris.Map{},
	}

	if err == nil {
		data = iris.Map{
			"user": a_user,
		}
	}

	return data
}
