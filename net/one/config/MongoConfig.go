package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var database string = "plant-api"

var c *mongo.Client

// 连接到MongoDB
func MongoDBConn() {
	url := "mongodb://admin:123456@47.96.7.148:27017"
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(url)
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	c = client
	log.Println("成功连接到MongoDB!")
}

func CollectionSave(key string, value string, classKey string, coll string) bool {
	// 指定获取要操作的数据集
	collection := c.Database(database).Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, bson.D{{"Info", value}, {classKey, key}})
	if err != nil {
		return false
	}
	//InsertedID := res.InsertedID
	return true
}

func CollectionGet(key string, data interface{}, classKey string, coll string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := c.Database(database).Collection(coll)
	one := collection.FindOne(ctx, bson.D{{classKey, key}})
	//r := entry.WeedResult{}
	err := one.Decode(data)
	if err != nil {
		log.Println("结果集", err)
		return false
	}
	//fmt.Println(r)
	return true
}

func Delete(key string, classKey string, coll string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := c.Database(database).Collection(coll)

	// context.TODO()
	deleteResult1, err := collection.DeleteOne(ctx, bson.D{{classKey, key}})
	if err != nil {
		log.Println(err)
		return false
	}
	log.Printf("Deleted %v documents in the trainers collection\n", deleteResult1.DeletedCount)
	return true

}

func DeleteAll(coll string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := c.Database(database).Collection(coll)

	// 删除所有context.TODO()
	deleteResult2, err := collection.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult2.DeletedCount)
	return true
}
