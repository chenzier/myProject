package common

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
)

//通用

// NewRedisClient 创建一个新的 Redis 客户端连接。
func NewRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456",
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

// String类型
// GetRedisString 通过键从 Redis 中获取一个字符串值。
func GetRedisString(ctx iris.Context, client *redis.Client, key string) (string, error) {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// SetRedisString 使用指定的键在 Redis 中设置一个字符串值。
func SetRedisString(ctx iris.Context, client *redis.Client, key string, value string) error {
	err := client.Set(ctx, key, value, 0).Err() //expiration表示过期时间
	if err != nil {
		return err
	}
	return nil
}

// DeleteRedisKey 从 Redis 中删除一个键。
func DeleteRedisKey(ctx iris.Context, client *redis.Client, key string) error {
	err := client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

//hash

// GetRedisHash 通过键从 Redis 中获取一个哈希值。
func GetRedisHash(ctx iris.Context, client *redis.Client, key string) (map[string]string, error) {
	val, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

// SetRedisHash 使用指定的键在 Redis 中设置一个哈希值。
func SetRedisHash(ctx iris.Context, client *redis.Client, key string, values map[string]interface{}) error {
	err := client.HMSet(ctx, key, values).Err()
	if err != nil {
		return err
	}
	return nil
}

// DeleteRedisHashField 从 Redis 中的哈希中删除一个字段。
func DeleteRedisHashField(ctx iris.Context, client *redis.Client, key string, field string) error {
	err := client.HDel(ctx, key, field).Err()
	if err != nil {
		return err
	}
	return nil
}

//List
//key ——> [1,2,3,...]
//
//func listFunc(){
//	cli,_:=
//}
