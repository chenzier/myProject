package main

import "fmt"

func main() {
	var a []int
	a = new([]int)
	a = append(a, 1)
	fmt.Println(a)
}

//package main
//
//import (
//	"context"
//	"fmt"
//	"github.com/go-redis/redis/v8"
//)
//
//func main() {
//	// 创建 Redis 客户端连接
//	client := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "123456",
//		DB:       0,
//	})
//
//	// 创建一个上下文对象
//	ctx := context.Background()
//
//	// 设置一个键值对
//	err := client.Set(ctx, "test_key", "test_value", 0).Err()
//	if err != nil {
//		fmt.Println("Failed to set key:", err)
//		return
//	}
//
//	// 查找键对应的值
//	val, err := client.Get(ctx, "test_key").Result()
//	if err != nil {
//		fmt.Println("Failed to get key:", err)
//		return
//	}
//
//	fmt.Println("Value of 'test_key':", val)
//}
//
////
////package main
////
////import (
////
////	"fmt"
////	//"bytes"
////	"net/http"
////	"net/http/httptest"
////
////)
////
////	func main() {
////		// 创建一个虚拟的 HTTP 请求
////		req := httptest.NewRequest("GET", "/product/all", nil)
////		// 创建一个 ResponseRecorder（它实现了 http.ResponseWriter 接口，用于捕获响应）
////		rr := httptest.NewRecorder()
////
////		// 这里模拟控制器的执行
////		ProductController_GetAll(rr, req)
////
////		// 检查响应状态码是否为 200
////		if status := rr.Code; status != http.StatusOK {
////			fmt.Printf("handler returned wrong status code: got %v want %v\n", status, http.StatusOK)
////		} else {
////			fmt.Println("handler returned correct status code:", status)
////		}
////
////		// 输出返回的页面内容
////		fmt.Println("Response body:", rr.Body.String())
////	}
////
////// 模拟控制器中的 GetAll 方法
////
////	func ProductController_GetAll(w http.ResponseWriter, r *http.Request) {
////		// 这里可以编写与 GetAll 方法相同的逻辑，例如获取产品信息并返回到页面
////		// 这里仅作为示例，返回一个简单的字符串
////		w.WriteHeader(http.StatusOK)
////		w.Write([]byte("This is a response from GetAll method"))
////	}
