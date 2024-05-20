package middlerware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisLock redis实现的分布式锁
type RedisLock struct {
	key        string
	value      string // 唯一标识,一般使用uuid
	expiration time.Duration
	redisCli   *redis.Client
	productID  int64 // 产品ID
}

func NewRedisLock(key, value string, expiration time.Duration, productID int64, cli *redis.Client) *RedisLock {
	if key == "" || value == "" || cli == nil {
		return nil
	}
	return &RedisLock{
		key:        key,
		value:      value,
		expiration: expiration,
		redisCli:   cli,
		productID:  productID,
	}
}

// Lock 添加分布式锁,
// rl.expiration过期时间,小于等于0,不过期
// 上锁后需要通过 UnLock方法释放锁
func (rl *RedisLock) Lock() (bool, error) {
	result, err := rl.redisCli.SetNX(context.Background(), rl.buildKey(), rl.value, rl.expiration).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

func (rl *RedisLock) TryLock(waitTime time.Duration) (bool, error) {
	var onceWaitTime = 20 * time.Millisecond
	if waitTime < onceWaitTime {
		waitTime = onceWaitTime
	}

	for index := 0; index < int(waitTime/onceWaitTime); index++ {
		locked, err := rl.Lock()
		if locked || err != nil {
			return locked, err
		}
		time.Sleep(onceWaitTime)
	}

	return false, nil
}

func (rl *RedisLock) UnLock() (bool, error) {
	script := redis.NewScript(`
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`)

	result, err := script.Run(context.Background(), rl.redisCli, []string{rl.buildKey()}, rl.value).Int64()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

// RefreshLock 存在则更新过期时间,不存在则创建key
// 首先检查指定键（KEYS[1]）是否存在
// 如果不存在，则使用 SETEX 命令创建该键，并设置过期时间为 ARGV[2]，同时将值设置为 ARGV[1]
// 然后返回 2 表示新创建了键。
// 如果键已经存在，脚本会检查当前键的值是否等于传入的参数值 ARGV[1]。
// 如果相等，说明当前锁是由当前客户端持有的，那么脚本会使用 EXPIRE 命令来刷新键的过期时间为 ARGV[2]。
// 如果不相等，说明当前锁不是由当前客户端持有的，那么脚本不会做任何操作，直接返回 0。
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

	result, err := script.Run(context.Background(), rl.redisCli, []string{rl.buildKey()}, rl.value, rl.expiration/time.Second).Int64()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

func (rl *RedisLock) buildKey() string {
	return fmt.Sprintf("%s:%d", rl.key, rl.productID)
}
