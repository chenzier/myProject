package main

import (
	"fmt"
	"log"
	"net/http"
	"product/fronted/middlerware"
	"sync"
	"time"
)

var sum int64 = 0

// 预存商品的数量
var productNum int64 = 1000000

// 互斥锁
var mutex sync.Mutex

// 获取秒杀商品
func GetOneProduct() bool {
	//加锁
	mutex.Lock()
	defer mutex.Unlock()
	//判断数据是否超限
	if sum < productNum {
		sum += 1
		return true
	}
	return false //商品抢购结束
}

func GetProduct(w http.ResponseWriter, req *http.Request) {
	fmt.Println(sum)
	if GetOneProduct() {
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
	return
}

func main4() {
	// 创建令牌桶限流器
	tbLimiter := middlerware.NewTokenBucketLimiter(10*time.Second, 1)

	http.HandleFunc("/getOne", middlerware.WrapWithLimiter(GetProduct, tbLimiter))

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatalf("Err:", err)
	}
}
