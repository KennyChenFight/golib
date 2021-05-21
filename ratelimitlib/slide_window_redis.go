package ratelimitlib

import (
	"context"
	"github.com/KennyChenFight/golib/redislib"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewSlideWindowRateLimiter(client *redislib.GORedisClient, capacity int64, interval time.Duration) *SlideWindowRateLimiter {
	script := redis.NewScript(SlideWindowIncrLuaScript)
	return &SlideWindowRateLimiter{client: client, luaScript: script, capacity: capacity, interval: interval.Nanoseconds(), expire: int64(interval.Seconds())}
}

type SlideWindowRateLimiter struct {
	client    *redislib.GORedisClient
	luaScript *redis.Script
	capacity  int64
	interval  int64
	expire    int64
}

const SlideWindowIncrLuaScript = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local interval = tonumber(ARGV[2])
local expire = tonumber(ARGV[3])
local lastTimestamp = tonumber(ARGV[4])

redis.call('ZREMRANGEBYSCORE', key, 0, lastTimestamp-interval)

local total = redis.call('ZCARD', key)
if total < capacity then
	redis.call('ZADD', key, lastTimestamp, lastTimestamp)
	total = total + 1
else
	total = -1
end

redis.call('EXPIRE', key, expire)

return total
`

func (s *SlideWindowRateLimiter) Incr(ctx context.Context, bucketName string, lastTimestamp int64) (int64, error) {
	args := []interface{}{s.capacity, s.interval, s.expire, lastTimestamp}

	total, err := s.luaScript.Run(ctx, s.client, []string{bucketName}, args...).Result()
	if err != nil {
		return 0, err
	}
	return total.(int64), nil
}
