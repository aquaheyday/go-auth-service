// internal/usecase/verify.go
// 이 파일은 이메일 인증 코드를 생성, 저장, 검증하는 비즈니스 로직(UseCase)을 정의합니다.
package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// VerificationRepository 인터페이스는 인증 코드의 저장, 확인, 조회, 삭제 기능을 정의합니다.
type VerificationRepository interface {
	SaveCode(ctx context.Context, email, code string) error           // 인증 코드를 저장 (만료 시간 포함)
	VerifyCode(ctx context.Context, email, code string) (bool, error) // 저장된 코드와 비교하여 일치 여부 반환
	GetCode(ctx context.Context, email string) (string, error)        // 원시 코드 조회 (추가 검증 시 사용)
	DeleteCode(ctx context.Context, email string) error               // 사용 후 코드 삭제
}

// MailSender 인터페이스는 이메일 전송 기능을 추상화합니다.
type MailSender interface {
	Send(to, subject, body string) error // 이메일 전송 구현체
}

// VerifyUseCase 인터페이스는 인증 코드 발송 및 검증 비즈니스 로직을 제공합니다.
type VerifyUseCase interface {
	SendVerification(ctx context.Context, email string) error         // 코드 생성 및 이메일 전송
	VerifyCode(ctx context.Context, email, code string) (bool, error) // 코드 검증
}

// verifyUseCase 구조체는 실제 레포지토리와 메일러를 사용하여 VerifyUseCase를 구현합니다.
type verifyUseCase struct {
	repo   VerificationRepository // 코드 저장소
	mailer MailSender             // 이메일 발송기
}

// NewVerifyUseCase 생성자 함수는 repo와 mailer를 주입받아 UseCase 인스턴스를 반환합니다.
func NewVerifyUseCase(repo VerificationRepository, mailer MailSender) VerifyUseCase {
	return &verifyUseCase{repo: repo, mailer: mailer}
}

// SendVerification은 랜덤 3바이트(6 hex 문자열) 코드를 생성하여 저장하고 이메일로 전송합니다.
func (v *verifyUseCase) SendVerification(ctx context.Context, email string) error {
	// 랜덤 바이트 생성
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	// 바이트를 hex 문자열로 인코딩
	code := hex.EncodeToString(b)

	// 레포지토리에 코드 저장 (예: Redis에 TTL 포함 저장)
	if err := v.repo.SaveCode(ctx, email, code); err != nil {
		return err
	}

	// 이메일 본문 생성 및 전송
	body := fmt.Sprintf("Your verification code is: %s", code)
	if err := v.mailer.Send(email, "Email Verification", body); err != nil {
		// 전송 실패 시, 필요하다면 저장된 코드 삭제 고려
		return err
	}

	return nil // 성공
}

// VerifyCode는 레포지토리를 통해 코드 일치 여부를 반환합니다.
func (v *verifyUseCase) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	// 저장된 코드와 비교
	ok, err := v.repo.VerifyCode(ctx, email, code)
	if err != nil {
		// 조회 에러 발생 시 전달
		return false, err
	}
	// 일치 여부 반환
	return ok, nil
}
