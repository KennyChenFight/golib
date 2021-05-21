// Package ratelimilib is for  rate limit any strategy operations
//
// As a quick start:
// 	client, err := redislib.NewGORedisClient(redislib.GORedisConfig{URL: "redis://localhost:6379"}, nil)
//	if err != nil {
//		panic(err)
//	}
//
//	limiter := ratelimitlib.NewSlideWindowRateLimiter(client, 5, 30 * time.Second)
// 	t := time.Now().UnixNano()
//	count, err := limiter.Incr(context.Background(), "kenny", t)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(count)
package ratelimitlib

import (
	"context"
)

type RateLimiter interface {
	Incr(ctx context.Context, bucketName string, lastTimestamp int64) (int64, error)
}
