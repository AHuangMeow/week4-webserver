package utils

import (
	"context"
	"fmt"
	"time"

	"week4-webserver/database"
)

type RedisBlacklist struct {
	ctx context.Context
}

func NewRedisBlacklist() *RedisBlacklist {
	return &RedisBlacklist{
		ctx: context.Background(),
	}
}

func (rb *RedisBlacklist) AddToken(token string, expireTime time.Time) error {
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	now := time.Now()
	ttl := expireTime.Sub(now)

	if ttl <= 0 {
		return fmt.Errorf("token is already expired")
	}

	err := database.RedisClient.Set(rb.ctx, getBlacklistKey(token), "blacklisted", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %v", err)
	}

	return nil
}

func (rb *RedisBlacklist) IsTokenBlacklisted(token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("token cannot be empty")
	}

	exists, err := database.RedisClient.Exists(rb.ctx, getBlacklistKey(token)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %v", err)
	}

	return exists > 0, nil
}

func (rb *RedisBlacklist) RemoveToken(token string) error {
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	err := database.RedisClient.Del(rb.ctx, getBlacklistKey(token)).Err()
	if err != nil {
		return fmt.Errorf("failed to remove token from blacklist: %v", err)
	}

	return nil
}

func (rb *RedisBlacklist) BatchAddTokens(tokens map[string]time.Time) error {
	if len(tokens) == 0 {
		return nil
	}

	pipe := database.RedisClient.Pipeline()

	for token, expireTime := range tokens {
		ttl := time.Until(expireTime)
		if ttl > 0 {
			pipe.Set(rb.ctx, getBlacklistKey(token), "blacklisted", ttl)
		}
	}

	_, err := pipe.Exec(rb.ctx)
	if err != nil {
		return fmt.Errorf("failed to batch add tokens to blacklist: %v", err)
	}

	return nil
}

func (rb *RedisBlacklist) CleanupExpiredTokens() (int64, error) {
	size, err := database.RedisClient.DBSize(rb.ctx).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get database size: %v", err)
	}

	return size, nil
}

func (rb *RedisBlacklist) GetBlacklistSize() (int64, error) {
	_, err := database.RedisClient.Info(rb.ctx, "keyspace").Result()
	if err != nil {
		return 0, err
	}

	size, err := database.RedisClient.DBSize(rb.ctx).Result()
	if err != nil {
		return 0, err
	}

	return size, nil
}

func getBlacklistKey(token string) string {
	return fmt.Sprintf("blacklist:token:%s", token)
}
