// internal/repository/token/token_repo.go

package token

import (
	"context"
	"time"
)

type Repository interface {
	StoreRefreshToken(ctx context.Context, userID, tokenID string, expiresAt time.Time) error
	ValidateRefreshToken(ctx context.Context, userID, tokenID string) (bool, error)
	DeleteRefreshToken(ctx context.Context, userID, tokenID string) error
	DeleteAllUserTokens(ctx context.Context, userID string) error
}
