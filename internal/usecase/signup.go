package usecase

import (
	"context"
	"errors"

	"github.com/aquaheyday/go-auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (string, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type SignupUsecase interface {
	SignUp(ctx context.Context, email, password, code string) (string, error)
}

type signupUsecase struct {
	userRepo         UserRepository
	verificationRepo VerificationRepository
}

func NewSignupUsecase(uRepo UserRepository, vRepo VerificationRepository) SignupUsecase {
	return &signupUsecase{userRepo: uRepo, verificationRepo: vRepo}
}

func (s *signupUsecase) SignUp(ctx context.Context, email, password, code string) (string, error) {
	valid, err := s.verificationRepo.VerifyCode(ctx, email, code)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", errors.New("invalid verification code")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashed),
	}
	return s.userRepo.Create(ctx, user)
}
