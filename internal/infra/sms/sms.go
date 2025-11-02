package sms

import "context"

// SMSProvider SMS 메시지 발송을 위한 인터페이스
type SMSProvider interface {
	// SendVerificationSMS 휴대폰 번호로 인증 코드를 발송합니다
	SendVerificationSMS(ctx context.Context, phoneNumber, code string) error
}
