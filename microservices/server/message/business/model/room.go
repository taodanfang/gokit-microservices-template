package model

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------

type Room struct {
	Room_uuid string   `json:"room_uuid", bson:"room_uuid"`
	Room_name string   `json:"room_name", bson:"room_name"`
	Members   []string `json:"members", bson:"members"`
}
