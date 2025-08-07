package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type VerificationRepository struct {
	rdb *redis.Client
}

func NewVerificationRepository(rdb *redis.Client) *VerificationRepository {
	return &VerificationRepository{rdb: rdb}
}

func (r *VerificationRepository) SaveCode(ctx context.Context, email, code string) error {
	key := "verify:" + email
	return r.rdb.Set(ctx, key, code, 10*time.Minute).Err()
}

func (r *VerificationRepository) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	key := "verify:" + email
	val, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return val == code, nil
}
