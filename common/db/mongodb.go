package db

import (
	"cctable/config"
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --------------------------------------------------------------------
// 全局变量
// --------------------------------------------------------------------

var (
	__mongo_db     *mongo.Database
	__mongo_client *mongo.Client
)

// --------------------------------------------------------------------
// 初始化方法
// --------------------------------------------------------------------

func Init_mongo_db() {
	cfg := config.Get_config()
	path := strings.Join(
		[]string{
			"mongodb://",
			cfg.Mongodb.Host,
			":",
			cfg.Mongodb.Port,
		}, "")

	log.Println("path:", path)

	client_Options := options.Client().ApplyURI(path)
	client, err := mongo.Connect(context.TODO(), client_Options)
	if err != nil {
		log.Fatalf("init mongodb (%s) is failure, %v\n", path, err)
		return
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(cfg.Mongodb.Database)

	__mongo_db = db
	__mongo_client = client

	log.Println("Mongodb connection is ok.")
}

// --------------------------------------------------------------------
// API：导出方法
// --------------------------------------------------------------------

func Get_mongo_db() *mongo.Database {
	return __mongo_db
}

func Get_mongo_client() *mongo.Client {
	return __mongo_client
}
