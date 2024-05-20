// middleware/limiter.go
package middlerware

import (
	"net/http"
	"time"

	"product/fronted/limiter"
)

func NewTokenBucketLimiter(rate time.Duration, capacity int64) *limiter.TokenBucketLimiter {
	return limiter.NewTokenBucketLimiter(rate, capacity)
}

// WrapWithLimiter 包装函数以应用限流器
func WrapWithLimiter(handler http.HandlerFunc, limiter *limiter.TokenBucketLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if limiter.Allow() {
			handler(w, req)
		} else {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		}
	}
}
