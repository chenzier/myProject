package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456",
		DB:       0,
	})
	redisLock := NewRedisLock("test_key1", "test_value1", 1*time.Minute, redisCli)
	//通过 NewRedisLock 函数创建了一个 RedisLock 实例 redisLock，用于管理分布式锁。
	locked, err := redisLock.TryLock(5 * time.Second) //接着，尝试获取锁并打印结果。
	fmt.Printf("locked:%v,err:%v\n", locked, err)
}

// RedisLock redis实现的分布式锁
type RedisLock struct {
	key   string //分布式锁在 Redis 中存储的key。每个锁都对应一个唯一的key
	value string //分布式锁的value，它是用来标识当前持有锁的唯一标识。在锁释放时，会根据这个值来验证是否可以释放。
	// value是唯一标识的,一般使用uuid
	expiration time.Duration //锁的过期时间，表示获取到的锁在 Redis 中的存储时长。如果设置为 0 或者一个负数，则表示锁不会自动过期，需要手动释放
	redisCli   *redis.Client // Redis 客户端连接对象的指针，用于与 Redis 服务器进行通信。通过这个客户端，可以执行获取锁、释放锁等操作
}

// redisLock := NewRedisLock("test_key1", "test_value1", 1*time.Minute, redisCli)
func NewRedisLock(key, value string, expiration time.Duration, cli *redis.Client) *RedisLock {
	if key == "" || value == "" || cli == nil {
		return nil
	}
	return &RedisLock{
		key:        key,
		value:      value,
		expiration: expiration,
		redisCli:   cli,
	}
}

// Lock 添加分布式锁,expiration过期时间,小于等于0,不过期,需要通过 UnLock方法释放锁
func (rl *RedisLock) Lock() (bool, error) { //返回的bool代表获取锁是否成功，error略
	result, err := rl.redisCli.SetNX(context.Background(), rl.key, rl.value, rl.expiration).Result()
	//调用了 Redis 客户端的 SetNX 方法，方法会检查key，如果不存在，将对应的value设为value
	//如果key已经存在，返回空值

	if err != nil {
		return false, err
	}

	return result, nil
}

// TryLock 方法是用来尝试获取锁的，如果获取失败，会在一定的等待时间内进行重试
func (rl *RedisLock) TryLock(waitTime time.Duration) (bool, error) {
	var onceWaitTime = 20 * time.Millisecond // 定义每次尝试获取锁的时间间隔，最小为20毫秒

	if waitTime < onceWaitTime { // 如果传入的等待时间小于等于每次尝试的间隔时间，将等待时间设置为每次尝试的间隔时间
		waitTime = onceWaitTime
	}

	//waitTime / onceWaitTime是总共需要尝试的次数
	for index := 0; index < int(waitTime/onceWaitTime); index++ {
		locked, err := rl.Lock()
		if locked || err != nil { // 如果成功获取锁或者出现了错误，则立即返回结果
			return locked, err
		}
		time.Sleep(onceWaitTime) // 如果未能获取锁，则等待一段时间再重试
	}

	return false, nil
}

func (rl *RedisLock) UnLock() (bool, error) {
	//// 创建 Lua 脚本
	script := redis.NewScript(`
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`)

	//// 运行 Lua 脚本，通过传入键和值进行匹配并删除键
	result, err := script.Run(context.Background(), rl.redisCli, []string{rl.key}, rl.value).Int64()
	if err != nil {
		return false, err
	}

	// 返回结果，如果成功删除了键，则返回 true，否则返回 false
	return result > 0, nil
}

// RefreshLock 存在则更新过期时间,不存在则创建key
func (rl *RedisLock) RefreshLock() (bool, error) {
	script := redis.NewScript(`
	local val = redis.call("GET", KEYS[1])
	if not val then
		redis.call("setex", KEYS[1], ARGV[2], ARGV[1])
		return 2
	elseif val == ARGV[1] then
		return redis.call("expire", KEYS[1], ARGV[2])
	else
		return 0
	end
	`)

	result, err := script.Run(context.Background(), rl.redisCli, []string{rl.key}, rl.value, rl.expiration/time.Second).Int64()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}
