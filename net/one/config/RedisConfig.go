package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"plant-api/net/one/entry"
	"time"
)

var (
	rdb *redis.Client
)

func RedisDB() *redis.Client {
	return rdb
}

// 初始化连接
func InitRedisClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "148.70.115.4:6379",
		Password: "MX5rlP9I62", // no password set
		DB:       1,            // use default DB
		PoolSize: 100,          // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

func SaveKV(key string, value string, exp time.Duration) error {
	ctx := context.Background()
	exp = exp * 1000 * 1000 * 1000
	err := rdb.Set(ctx, key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func SaveKJson(key string, value interface{}, exp time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	valString := string(b)
	err = SaveKV(key, valString, exp)
	if err != nil {
		return err
	}
	return nil
}

// 获取
func GetK(key string) (string, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "nil", err
	}
	return val, nil
}

// 获取
func GetJson(key string, typeI interface{}) error {
	val, err := GetK(key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), typeI)
	if err != nil {
		return err
	}
	return nil
}

// DEP
func SaveEntry() {
	ctx := context.Background()
	if err := InitRedisClient(); err != nil {
		fmt.Println("创建链接失败")
		return
	}

	user := entry.SysUser{
		Id:       0,
		UUID:     "didididid",
		OpenId:   "1233213opid",
		UserName: "tom",
		Province: "jfasofj",
	}

	key := "test-user-1"
	val := make(map[string]interface{})
	val["u1"] = user

	_, err := rdb.HSetNX(ctx, key, "u1", user).Result()
	if err == nil {
		fmt.Println("ok")
	} else {
		fmt.Println("err", err.Error())
	}

	//b, err := rdb.HMSet(ctx, key, val).Result()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//if !b {
	//	fmt.Println("err save")
	//	return
	//}
	//
	//i, err := rdb.HMGet(ctx, key).Result()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(i)

}

func V8Example() {
	ctx := context.Background()
	if err := InitRedisClient(); err != nil {
		fmt.Println("创建链接失败")
		return
	}

	key := "lwm-test-key"
	key2 := "lwm-test2-key"
	value := "hello world"

	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		fmt.Println("保存key失败")
		fmt.Println(err)
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key[lwm-test-key] ： ", val)

	val2, err := rdb.Get(ctx, key2).Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
