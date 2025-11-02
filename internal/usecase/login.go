// internal/usecase/login.go

package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/aquaheyday/go-auth-service/internal/repository/postgres"
	tokenRepo "github.com/aquaheyday/go-auth-service/internal/repository/token"
	"github.com/aquaheyday/go-auth-service/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase interface {
	Login(ctx context.Context, email, password string) (string, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
}

type loginUseCase struct {
	userRepo  postgres.UserRepository
	tokenRepo tokenRepo.Repository
}

func NewLoginUseCase(userRepo postgres.UserRepository, tokenRepo tokenRepo.Repository) LoginUseCase {
	return &loginUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

// 로그인 처리
func (uc *loginUseCase) Login(ctx context.Context, email, password string) (string, string, string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", "", err
	}

	// 비밀번호 확인
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", "", errors.New("invalid credentials")
	}

	// 액세스 토큰 생성
	accessToken, err := token.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", "", err
	}

	// 리프레시 토큰 생성
	refreshToken, tokenID, err := token.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", "", err
	}

	// 리프레시 토큰을 Redis에 저장
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := uc.tokenRepo.StoreRefreshToken(ctx, user.ID, tokenID, expiresAt); err != nil {
		return "", "", "", err
	}

	return user.ID, accessToken, refreshToken, nil
}

// 토큰 갱신
func (uc *loginUseCase) RefreshToken(ctx context.Context, refreshTokenStr string) (string, string, error) {
	// 리프레시 토큰 검증
	claims, err := token.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return "", "", err
	}

	// Redis에서 토큰 확인
	valid, err := uc.tokenRepo.ValidateRefreshToken(ctx, claims.UserID, claims.TokenID)
	if err != nil {
		return "", "", err
	}
	if !valid {
		return "", "", errors.New("invalid or expired refresh token")
	}

	// 기존 토큰 삭제
	if err := uc.tokenRepo.DeleteRefreshToken(ctx, claims.UserID, claims.TokenID); err != nil {
		return "", "", err
	}

	// 새 액세스 토큰 생성
	accessToken, err := token.GenerateAccessToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	// 새 리프레시 토큰 생성 (토큰 회전)
	newRefreshToken, newTokenID, err := token.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	// 새 리프레시 토큰 저장
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := uc.tokenRepo.StoreRefreshToken(ctx, claims.UserID, newTokenID, expiresAt); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// 로그아웃
func (uc *loginUseCase) Logout(ctx context.Context, refreshTokenStr string) error {
	// 리프레시 토큰 검증
	claims, err := token.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return err
	}

	// Redis에서 토큰 삭제
	return uc.tokenRepo.DeleteRefreshToken(ctx, claims.UserID, claims.TokenID)
}
