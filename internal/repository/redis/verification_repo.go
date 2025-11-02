package redis

import (
	"context"
	"fmt"
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

// GetCode retrieves the raw verification code.
func (r *VerificationRepository) GetCode(ctx context.Context, email string) (string, error) {
	key := "verify:" + email
	return r.rdb.Get(ctx, key).Result()
}

// DeleteCode removes the verification code after use.
func (r *VerificationRepository) DeleteCode(ctx context.Context, email string) error {
	key := "verify:" + email
	return r.rdb.Del(ctx, key).Err()
}

// StorePhoneVerificationCode 휴대폰 인증 코드를 저장합니다
func (r *VerificationRepository) StorePhoneVerificationCode(ctx context.Context, phoneNumber, code string, expiration time.Duration) error {
	key := fmt.Sprintf("phone_verification:%s", phoneNumber)
	return r.client.Set(ctx, key, code, expiration)
}

// GetPhoneVerificationCode 저장된 휴대폰 인증 코드를 조회합니다
func (r *verificationRepository) GetPhoneVerificationCode(ctx context.Context, phoneNumber string) (string, error) {
	key := fmt.Sprintf("phone_verification:%s", phoneNumber)
	return r.client.Get(ctx, key)
}

// DeletePhoneVerificationCode 휴대폰 인증 코드를 삭제합니다
func (r *verificationRepository) DeletePhoneVerificationCode(ctx context.Context, phoneNumber string) error {
	key := fmt.Sprintf("phone_verification:%s", phoneNumber)
	return r.client.Del(ctx, key)
}
