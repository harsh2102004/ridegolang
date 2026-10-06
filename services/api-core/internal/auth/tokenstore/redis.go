package tokenstore

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var RDB *redis.Client

// Initialize Redis (call this in main)
func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

// Store token → userID 
func StoreVerificationToken(token string, userID string, ttl time.Duration) error {
	key := "email_verify:" + token
	return RDB.Set(ctx, key, userID, ttl).Err()
}

// Get userID from token
func GetUserIDByToken(token string) (string, error) {
	key := "email_verify:" + token
	return RDB.Get(ctx, key).Result()
}

// Delete token after use
func DeleteToken(token string) error {
	key := "email_verify:" + token
	return RDB.Del(ctx, key).Err()
}

func StoreRefreshToken(ctx context.Context, token string, userID string) error {
	key := "refresh:" + token
	return RDB.Set(ctx, key, userID, 7*24*time.Hour).Err()
}

func GetRefreshToken(ctx context.Context, token string) (string, error) {
	key := "refresh:" + token
	return RDB.Get(ctx, key).Result()
}

func DeleteRefreshToken(ctx context.Context, token string) error {
	key := "refresh:" + token
	return RDB.Del(ctx, key).Err()
}