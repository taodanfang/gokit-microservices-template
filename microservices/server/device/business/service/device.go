package service

import (
	"cctable/common/results"
	"cctable/common/tools"
	"cctable/microservices/server/device/business/model"
	"context"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type I_device_service interface {
	New_device(ctx context.Context, device_name, device_code string) *results.Result

	Update_device(ctx context.Context, device_uuid string, item_key string, item_value interface{}) *results.Result

	Check_device_by_name_and_code(ctx context.Context, device_name, device_code string) *results.Result
	Get_device_by_uuid(ctx context.Context, device_uuid string) *results.Result
	Get_all_devices(ctx context.Context) *results.Result
}

type Device_service struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// --------------------------------------------------------------------
// 构造函数
// --------------------------------------------------------------------

func New_device_service(db *mongo.Database) *Device_service {
	service := Device_service{}
	service.db = db
	service.collection = db.Collection("devices")
	return &service
}

// --------------------------------------------------------------------
// 接口方法
// --------------------------------------------------------------------

func (s *Device_service) New_device(ctx context.Context, device_name, device_code string) *results.Result {
	data := results.API()

	a_device := model.Device{
		Device_uuid: primitive.NewObjectID().Hex(),
		Device_name: device_name,
		Device_code: device_code,
		Create_at:   time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, a_device)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_insert.Error())
	}

	current := s.get_device(ctx, a_device.Device_uuid)
	data["device"] = current["device"]

	return results.Ok(data)
}

func (s *Device_service) Update_device(ctx context.Context, device_uuid string, item_key string, item_value interface{}) *results.Result {
	data := results.API()

	filter := bson.D{
		{"device_uuid", device_uuid},
	}

	update := bson.D{{}}

	switch item_key {
	case "device_name":
		update = bson.D{{"$set", bson.D{{item_key, item_value.(string)}}}}
		break
	case "device_code":
		update = bson.D{{"$set", bson.D{{item_key, item_value.(string)}}}}
		break
	default:
		return results.Error(data, results.Err_db_invalid_item_key_for_update.Error())
	}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		tools.Log(err)
		return results.Error(data, results.Err_db_failed_to_update.Error())
	}

	current := s.get_device(ctx, device_uuid)
	data["device"] = current["device"]

	return results.Ok(data)
}

func (s *Device_service) Get_device_by_uuid(ctx context.Context, device_uuid string) *results.Result {
	data := results.API()

	filter := bson.D{
		{"device_uuid", device_uuid},
	}

	var a_device = model.Device{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_device)

	if err != nil {
		return results.Error(data, results.Err_device_does_not_exist.Error())
	}

	data["device"] = a_device

	return results.Ok(data)
}

func (s *Device_service) Check_device_by_name_and_code(ctx context.Context, device_name, device_code string) *results.Result {
	data := results.API()

	filter := bson.D{
		{"device_name", device_name},
		{"device_code", device_code},
	}

	var a_device = model.Device{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_device)

	if err != nil {
		return results.Error(data, results.Err_device_does_not_exist.Error())
	}

	data["device"] = a_device

	return results.Ok(data)
}

func (s *Device_service) Get_all_devices(ctx context.Context) *results.Result {
	data := results.API()

	filter := bson.D{{}}

	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		return results.Error(data, results.Err_db_failed_to_find.Error())
	}

	var devices []interface{}
	for cur.Next(ctx) {
		var a_device model.Device
		err := cur.Decode(&a_device)
		if err != nil {
			return results.Error(data, results.Err_db_failed_to_find.Error())
		}
		devices = append(devices, a_device)
	}

	_ = cur.Close(ctx)

	if len(devices) <= 0 {
		data["devices"] = tools.EmptyArrayString
	} else {
		data["devices"] = devices
	}

	return results.Ok(data)
}

// --------------------------------------------------------------------
// 辅助方法
// --------------------------------------------------------------------

func (s *Device_service) get_device(ctx context.Context, device_uuid string) (data iris.Map) {

	filter := bson.D{{"device_uuid", device_uuid}}

	var a_device = model.Device{}
	err := s.collection.FindOne(ctx, filter).Decode(&a_device)

	data = iris.Map{
		"device": iris.Map{},
	}

	if err == nil {
		data = iris.Map{
			"device": a_device,
		}
	}

	return data
}
