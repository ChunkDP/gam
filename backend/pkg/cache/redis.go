package cache

import (
	"context"
	"encoding/json"
	"normaladmin/backend/config"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	cfg    config.RedisConfig
)

// InitRedis 初始化Redis连接
func InitRedis(redisConfig config.RedisConfig) error {
	cfg = redisConfig
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return client.Ping(context.Background()).Err()
}

// Get 获取缓存
func Get(key string) (string, error) {
	return client.Get(context.Background(), key).Result()
}

// GetObject 获取并解析JSON对象
func GetObject(key string, val interface{}) error {
	data, err := Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), val)
}

// Set 设置缓存
func Set(key string, value string, expiration ...time.Duration) error {
	exp := time.Duration(cfg.DefaultTTL) * time.Second
	//24 * time.Hour // 默认24小时
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	return client.Set(context.Background(), key, value, exp).Err()
}

// SetObject 设置JSON对象
func SetObject(key string, val interface{}, expiration ...time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return Set(key, string(data), expiration...)
}

// SetNX 如果key不存在则设置
func SetNX(key string, value interface{}, expiration ...time.Duration) (bool, error) {
	exp := time.Duration(cfg.DefaultTTL) * time.Second
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	return client.SetNX(context.Background(), key, value, exp).Result()
}

// Delete 删除缓存
func Delete(key string) error {
	return client.Del(context.Background(), key).Err()
}

// Exists 检查key是否存在
func Exists(key string) bool {
	result, _ := client.Exists(context.Background(), key).Result()
	return result > 0
}

// Lock 获取分布式锁
func Lock(key string, expiration ...time.Duration) bool {
	exp := time.Duration(cfg.LockTimeout) * time.Second
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	success, _ := SetNX("lock:"+key, 1, exp)
	return success
}

// Unlock 释放分布式锁
func Unlock(key string) error {
	return Delete("lock:" + key)
}

// HSet 设置哈希表字段
func HSet(key, field string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return client.HSet(context.Background(), key, field, data).Err()
}

// HGet 获取哈希表字段
func HGet(key, field string, val interface{}) error {
	data, err := client.HGet(context.Background(), key, field).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), val)
}

// HDel 删除哈希表字段
func HDel(key string, fields ...string) error {
	return client.HDel(context.Background(), key, fields...).Err()
}

// Incr 递增
func Incr(key string) (int64, error) {
	return client.Incr(context.Background(), key).Result()
}

// Decr 递减
func Decr(key string) (int64, error) {
	return client.Decr(context.Background(), key).Result()
}

// Close 关闭Redis连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// DeletePattern 删除匹配模式的所有键
func DeletePattern(pattern string) error {
	ctx := context.Background()
	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
