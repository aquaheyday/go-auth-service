package sms

import (
	"context"
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/aquaheyday/go-auth-service/pkg/config"
)

type twilioProvider struct {
	client *twilio.RestClient
	from   string
}

// NewTwilioProvider Twilio SMS 서비스 구현체를 생성합니다
func NewTwilioProvider(config config.TwilioConfig) SMSProvider {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	return &twilioProvider{
		client: client,
		from:   config.FromNumber,
	}
}

// SendVerificationSMS Twilio API를 사용해 SMS 인증 코드를 발송합니다
func (p *twilioProvider) SendVerificationSMS(ctx context.Context, phoneNumber, code string) error {
	params := &twilioApi.CreateMessageParams{
		To:   &phoneNumber,
		From: &p.from,
		Body: fmt.Sprintf("인증 코드: %s", code),
	}

	_, err := p.client.Api.CreateMessage(params)
	return err
}
