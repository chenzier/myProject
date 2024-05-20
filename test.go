package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"product/common"
)

type User struct {
	Username string `json:"username"`
}

func (u *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}
func main() {
	cli, err := common.NewRedisClient()
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	res := cli.Set(ctx, "user", &User{Username: "user1"}, 0)
	if res.Err() != nil {
		log.Fatal(res.Err())
	}
	getRes, err := cli.Get(ctx, "user").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(getRes)

	//cli.Set(ctx, "exm", "exm after 3s", 3*time.Second)
	//tk := time.NewTicker(1 * time.Second)
	//for range tk.C {
	//	expres, err := cli.Get(ctx, "exm").Result()
	//	if err != nil {
	//		fmt.Println("找不到缓存")
	//		break
	//	}
	//	fmt.Println(expres)
	//}

	cli.LPush(ctx, "list1", "v1", "v2", "v3")   //注意是Left Push，先进后出 即List内是v3,v2,v1
	cli.RPush(ctx, "list1", "v4", "v5", "v6")   //双端队列的Right Push
	getRes, _ = cli.LPop(ctx, "list1").Result() //pop一个值出来
	fmt.Println(getRes)
}

//import (
//	"fmt"
//	"sync"
//	"time"
//)
//
//type TokenBucket struct {
//	rate       int
//	capacity   int
//	tokens     int
//	lastRefill time.Time
//	mu         sync.Mutex // 互斥锁
//}
//
//func NewTokenBucket(rate, capacity int) *TokenBucket {
//	tb := &TokenBucket{
//		rate:       rate,
//		capacity:   capacity,
//		tokens:     capacity,
//		lastRefill: time.Now(),
//	}
//	return tb
//}
//
//func (tb *TokenBucket) Take() bool {
//	tb.mu.Lock()
//	defer tb.mu.Unlock()
//
//	now := time.Now()
//	elapsed := now.Sub(tb.lastRefill)
//
//	refillCount := int(elapsed.Seconds()) * tb.rate
//	if refillCount > 0 {
//		tb.tokens = tb.tokens + refillCount
//		if tb.tokens > tb.capacity {
//			tb.tokens = tb.capacity
//		}
//		tb.lastRefill = now
//	}
//
//	if tb.tokens > 0 {
//		tb.tokens--
//		return true
//	}
//
//	return false
//}
//
//func main() {
//	tb := NewTokenBucket(5, 10)
//
//	// 使用 WaitGroup 来等待所有 goroutine 完成
//	var wg sync.WaitGroup
//	for i := 1; i <= 100; i++ {
//		wg.Add(1)
//		go func(i int) {
//			defer wg.Done()
//			if tb.Take() {
//				fmt.Println("Request", i, "is processed.")
//			} else {
//				fmt.Println("Request", i, "is rejected.")
//			}
//		}(i)
//	}
//
//	// 等待所有 goroutine 完成
//	wg.Wait()
//}
