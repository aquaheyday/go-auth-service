package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type VerificationRepository interface {
	SaveCode(ctx context.Context, email, code string) error
	VerifyCode(ctx context.Context, email, code string) (bool, error)
}

type MailSender interface {
	Send(to, subject, body string) error
}

type VerifyUseCase interface {
	SendVerification(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email, code string) (bool, error)
}

type verifyUseCase struct {
	repo   VerificationRepository
	mailer MailSender
}

func NewVerifyUseCase(repo VerificationRepository, mailer MailSender) VerifyUseCase {
	return &verifyUseCase{repo: repo, mailer: mailer}
}

func (v *verifyUseCase) SendVerification(ctx context.Context, email string) error {
	// 랜덤 코드 생성
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	code := hex.EncodeToString(b)

	// Redis에 저장
	if err := v.repo.SaveCode(ctx, email, code); err != nil {
		return err
	}

	// 이메일 전송
	body := fmt.Sprintf("Your verification code is: %s", code)
	return v.mailer.Send(email, "Email Verification", body)
}

func (v *verifyUseCase) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	return v.repo.VerifyCode(ctx, email, code)
}
