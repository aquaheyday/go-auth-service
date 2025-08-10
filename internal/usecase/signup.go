// internal/usecase/signup.go
// 이 파일은 회원 가입(SignUp) 비즈니스 로직을 정의합니다.
package usecase

import (
	"context"
	"errors"

	"github.com/aquaheyday/go-auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository 인터페이스는 사용자 생성 및 조회 기능을 추상화합니다.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (string, error)      // 새 사용자 생성 및 ID 반환
	GetByEmail(ctx context.Context, email string) (*domain.User, error) // 이메일로 사용자 조회
}

// SignupUseCase 인터페이스는 회원 가입 흐름(SignUp)을 정의합니다.
type SignupUseCase interface {
	SignUp(ctx context.Context, email, password, code string) (string, error) // 이메일, 비밀번호, 인증 코드로 가입 처리
}

// signupUseCase 구조체는 UserRepository와 VerificationRepository를 사용하여 SignUp 로직을 구현합니다.
type signupUseCase struct {
	userRepo         UserRepository         // 사용자 저장소
	verificationRepo VerificationRepository // 인증 코드 검증 저장소
}

// NewSignupUseCase 생성자 함수는 필요한 저장소를 주입받아 SignupUseCase 인스턴스를 반환합니다.
func NewSignupUseCase(uRepo UserRepository, vRepo VerificationRepository) SignupUseCase {
	return &signupUseCase{userRepo: uRepo, verificationRepo: vRepo}
}

// SignUp 메서드는 다음 순서로 회원 가입을 처리합니다:
// 1) 인증 코드 검증
// 2) 비밀번호 해시 생성
// 3) 사용자 도메인 모델 생성 및 저장
// 4) 생성된 사용자 ID 반환
func (s *signupUseCase) SignUp(ctx context.Context, email, password, code string) (string, error) {
	// 인증 코드 검증
	valid, err := s.verificationRepo.VerifyCode(ctx, email, code)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", errors.New("invalid verification code")
	}

	// 비밀번호를 bcrypt로 해시 처리
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 도메인 모델 생성
	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashed),
	}

	// 저장소에 사용자 생성 요청 및 ID 반환
	return s.userRepo.Create(ctx, user)
}
