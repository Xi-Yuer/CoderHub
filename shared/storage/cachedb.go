package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig Redis配置结构
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int           // 连接池大小
	MinIdleConns int           // 最小空闲连接数
	DefaultTTL   time.Duration // 默认过期时间
}

// RedisDB 接口定义
type RedisDB interface {
	Get(key string) (string, error)
	GetWithContext(ctx context.Context, key string) (string, error)
	Set(key string, value string) error
	HSet(key string, field string, value interface{}) error
	HMGet(key string, fields ...string) ([]interface{}, error)
	SetWithTTL(key string, value string, expiration time.Duration) error
	SetNX(key string, value string, expiration time.Duration) (bool, error)
	Del(key ...string) error
	Expire(key string, expiration time.Duration) error
	Exists(key string) (bool, error)
	NewScript(script string) *redis.Script
	Close() error
	Pipeline() redis.Pipeliner
	HIncrBy(key string, field string, value int64) (int64, error)
	HGet(key string, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	SMembers(key string) ([]string, error)
	SAdd(key string, members ...interface{}) error
	SRem(key string, members ...interface{}) error
	LRange(key string, start, stop int64) ([]string, error)
	RPush(key string, values ...interface{}) error
	LRem(key string, count int64, value interface{}) error
}

// DefaultConfig 默认配置
func DefaultConfig() *RedisConfig {
	return &RedisConfig{
		Addr:         "8.156.66.210:6379",
		Password:     "2214380963Wx!!",
		DB:           0,
		PoolSize:     50,             // 增加连接池大小
		MinIdleConns: 10,             // 增加最小空闲连接数
		DefaultTTL:   24 * time.Hour, // 默认24小时过期
	}
}

// NewRedisDB 创建Redis客户端
func NewRedisDB(config *RedisConfig) (RedisDB, error) {
	if config == nil {
		config = DefaultConfig()
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis连接失败: %v", err)
	}

	return &RedisDBImpl{
		Client:     client,
		DefaultTTL: config.DefaultTTL,
	}, nil
}

type RedisDBImpl struct {
	Client     *redis.Client
	DefaultTTL time.Duration
}

func (r *RedisDBImpl) Get(key string) (string, error) {
	return r.GetWithContext(context.Background(), key)
}

func (r *RedisDBImpl) GetWithContext(ctx context.Context, key string) (string, error) {
	result, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("key不存在: %s", key)
	}
	return result, err
}

func (r *RedisDBImpl) Set(key string, value string) error {
	return r.SetWithTTL(key, value, r.DefaultTTL)
}

func (r *RedisDBImpl) HSet(key string, field string, value interface{}) error {
	ctx := context.Background()
	return r.Client.HSet(ctx, key, field, value).Err()
}

func (r *RedisDBImpl) HMGet(key string, fields ...string) ([]interface{}, error) {
	ctx := context.Background()
	return r.Client.HMGet(ctx, key, fields...).Result()
}

func (r *RedisDBImpl) SetWithTTL(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisDBImpl) SetNX(key string, value string, expiration time.Duration) (bool, error) {
	ctx := context.Background()
	return r.Client.SetNX(ctx, key, value, expiration).Result()
}

func (r *RedisDBImpl) Del(keys ...string) error {
	ctx := context.Background()
	return r.Client.Del(ctx, keys...).Err()
}

func (r *RedisDBImpl) Expire(key string, expiration time.Duration) error {
	ctx := context.Background()
	return r.Client.Expire(ctx, key, expiration).Err()
}

func (r *RedisDBImpl) Exists(key string) (bool, error) {
	ctx := context.Background()
	n, err := r.Client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *RedisDBImpl) Pipeline() redis.Pipeliner {
	return r.Client.Pipeline()
}

func (r *RedisDBImpl) NewScript(script string) *redis.Script {
	return redis.NewScript(script)
}

func (r *RedisDBImpl) Close() error {
	return r.Client.Close()
}

func (r *RedisDBImpl) HIncrBy(key string, field string, value int64) (int64, error) {
	ctx := context.Background()
	return r.Client.HIncrBy(ctx, key, field, value).Result()
}

func (r *RedisDBImpl) HGet(key string, field string) (string, error) {
	ctx := context.Background()
	return r.Client.HGet(ctx, key, field).Result()
}

func (r *RedisDBImpl) HGetAll(key string) (map[string]string, error) {
	ctx := context.Background()
	return r.Client.HGetAll(ctx, key).Result()
}

func (r *RedisDBImpl) SMembers(key string) ([]string, error) {
	ctx := context.Background()
	return r.Client.SMembers(ctx, key).Result()
}

func (r *RedisDBImpl) SAdd(key string, members ...interface{}) error {
	return r.Client.SAdd(context.Background(), key, members...).Err()
}

func (r *RedisDBImpl) SRem(key string, members ...interface{}) error {
	return r.Client.SRem(context.Background(), key, members...).Err()
}

func (r *RedisDBImpl) LRange(key string, start, stop int64) ([]string, error) {
	ctx := context.Background()
	return r.Client.LRange(ctx, key, start, stop).Result()
}

func (r *RedisDBImpl) RPush(key string, values ...interface{}) error {
	ctx := context.Background()
	return r.Client.RPush(ctx, key, values...).Err()
}
func (r *RedisDBImpl) LRem(key string, count int64, value interface{}) error {
	ctx := context.Background()
	return r.Client.LRem(ctx, key, count, value).Err()
}
