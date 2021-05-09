// Package redislib is for encapsulating github.com/go-redis/redis any operations
//
// As a quick start:
// 	client, err := redislib.NewGORedisClient(redislib.GORedisConfig{URL: "redis://localhost:6379"}, nil)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(client.Ping(context.Background()).Err())
//
//	client, err := redislib.NewGORedisClusterClient(redislib.GORedisClusterConfig{Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"}}, nil)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(client.Ping(context.Background()).Err())
package redislib

import (
	"github.com/go-redis/redis/v8"
)

func NewGORedisClient(config GORedisConfig, customOpts *redis.Options) (*GORedisClient, error) {
	var opts redis.Options
	defaultOpts, err := redis.ParseURL(config.URL)
	if err != nil {
		return nil, err
	}
	opts = *defaultOpts
	if customOpts != nil {
		opts = *customOpts
	}
	return &GORedisClient{redis.NewClient(&opts)}, nil
}

func NewGORedisClusterClient(config GORedisClusterConfig, customOpts *redis.ClusterOptions) (*GORedisClusterClient, error) {
	var opts redis.ClusterOptions
	opts.Addrs = config.Addrs
	if customOpts != nil {
		opts = *customOpts
	}
	return &GORedisClusterClient{redis.NewClusterClient(&opts)}, nil
}

type GORedisConfig struct {
	URL string
}

type GORedisClusterConfig struct {
	Addrs []string
}

type GORedisClient struct {
	*redis.Client
}

type GORedisClusterClient struct {
	*redis.ClusterClient
}
