package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/aquaheyday/go-auth-service/pkg/token"
)

// 인증 흐름 전체를 담당하는 UseCase
type VerifyAuthUseCase interface {
	// 코드 검증 후, 기존 회원이면 토큰 반환, 신규면 가입 필요 표시
	VerifyAndAuth(ctx context.Context, email, code string) (token string, signUpRequired bool, err error)
}

type verifyAuthUseCase struct {
	repo      VerificationRepository
	userRepo  UserRepository
	tokenProv token.Provider
}

func NewVerifyAuthUseCase(
	repo VerificationRepository,
	userRepo UserRepository,
	tp token.Provider,
) VerifyAuthUseCase {
	return &verifyAuthUseCase{repo: repo, userRepo: userRepo, tokenProv: tp}
}

func (u *verifyAuthUseCase) VerifyAndAuth(ctx context.Context, email, code string) (string, bool, error) {
	// 1) 코드 확인 (GetCode, compare, DeleteCode)
	stored, err := u.repo.GetCode(ctx, email)
	if err != nil {
		return "", false, fmt.Errorf("code missing: %w", err)
	}
	if stored != code {
		return "", false, errors.New("code mismatch")
	}
	if err := u.repo.DeleteCode(ctx, email); err != nil {
		return "", false, err
	}

	// 2) 사용자 조회
	userID, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// not found → 신규 가입 필요
		return "", true, nil
	}

	// 3) 기존 사용자 → JWT 발급
	tok, err := u.tokenProv.GenerateAccessToken(fmt.Sprint(userID))
	if err != nil {
		return "", false, err
	}
	return tok, false, nil
}
