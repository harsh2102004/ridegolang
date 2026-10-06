package tokenstore

import (
	"context"
	"fmt"
	"time"
)

func StorePasswordResetToken(
	ctx context.Context,
	token string,
	userId string,
	ttl time.Duration,
) error {
	key := "pwd_reset:" + token
	return RDB.Set(ctx, key, userId, ttl).Err()
}

func GetPasswordResetToken(
	ctx context.Context,
	token string,
) (string, error) {
	key := "pwd_reset:" + token
    fmt.Println("REDIS KEY:", key)
	return RDB.Get(ctx, key).Result()
}

func DeletePasswordResetToken(
	ctx context.Context,
	token string,
) error {
	key := "pwd_reset" + token
	return RDB.Del(ctx, key).Err()
}
