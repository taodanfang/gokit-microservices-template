package service

import (
	"cctable/common/results"
	"cctable/microservices/server/message/business/model"
	"context"
	"github.com/kataras/iris/v12"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_room_service interface {
	New_room(ctx context.Context, room_name string) *results.Result
	Get_room_by_uuid(ctx context.Context, room_uuid string) *results.Result
	Check_room_by_name(ctx context.Context, room_name string) *results.Result
}

type Room_service struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_room_service(db *mongo.Database) *Room_service {
	service := Room_service{}
	service.db = db
	service.collection = db.Collection("rooms")
	return &service
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *Room_service) New_room(ctx context.Context, room_name string) *results.Result {
	data := results.API()

	a_room := model.Room{
		Room_uuid: primitive.NewObjectID().Hex(),
		Room_name: room_name,
	}

	_, err := s.collection.InsertOne(ctx, a_room)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_insert.Error())
	}

	current := s.get_room(ctx, a_room.Room_uuid)
	data["room"] = current["room"]

	return results.Ok(data)
}

func (s *Room_service) Get_room_by_uuid(ctx context.Context, room_uuid string) *results.Result {
	data := results.API()

	filter := bson.D{
		{"room_uuid", room_uuid},
	}

	var a_room = model.Room{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_room)

	if err != nil {
		return results.Error(data, results.Err_room_does_not_exist.Error())
	}

	data["room"] = a_room

	return results.Ok(data)
}

func (s *Room_service) Check_room_by_name(ctx context.Context, room_name string) *results.Result {
	data := results.API()

	filter := bson.D{
		{"room_name", room_name},
	}

	var a_room = model.Room{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_room)

	if err != nil {
		return results.Error(data, results.Err_room_does_not_exist.Error())
	}

	data["room"] = a_room

	return results.Ok(data)
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (s *Room_service) get_room(ctx context.Context, room_uuid string) (data iris.Map) {

	filter := bson.D{{"room_uuid", room_uuid}}

	var a_room = model.Room{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_room)

	data = iris.Map{
		"room": iris.Map{},
	}

	if err == nil {
		data = iris.Map{
			"room": a_room,
		}
	}

	return data
}
