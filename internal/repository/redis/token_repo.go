// internal/repository/redis/token_repo.go

package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type tokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) *tokenRepository {
	return &tokenRepository{
		client: client,
	}
}

func (r *tokenRepository) StoreRefreshToken(ctx context.Context, userID, tokenID string, expiresAt time.Time) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, tokenID)
	duration := time.Until(expiresAt)

	return r.client.Set(ctx, key, "valid", duration).Err()
}

func (r *tokenRepository) ValidateRefreshToken(ctx context.Context, userID, tokenID string) (bool, error) {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, tokenID)

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return val == "valid", nil
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, userID, tokenID string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, tokenID)

	return r.client.Del(ctx, key).Err()
}

func (r *tokenRepository) DeleteAllUserTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("refresh_token:%s:*", userID)

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}
