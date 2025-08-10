package mailer

import (
	"fmt"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	apiKey    string
	fromEmail string
	fromName  string
	// 개발용 샌드박스 모드: 실제 발송 대신 검증만 수행
	sandbox bool
}

func NewSendGridMailer(apiKey, fromEmail, fromName string, sandbox bool) *SendGridMailer {
	return &SendGridMailer{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
		sandbox:   sandbox,
	}
}

func boolPtr(b bool) *bool { return &b }

// UseCase에서 요구하는 시그니처: 텍스트 바디 하나만 받음
// 여기서 텍스트를 기반으로 HTML 바디를 자동으로 만들어 같이 보냅니다.
func (m *SendGridMailer) Send(to, subject, body string) error {
	if m.apiKey == "" {
		m.apiKey = os.Getenv("SENDGRID_API_KEY")
	}
	if m.fromEmail == "" {
		m.fromEmail = os.Getenv("SENDGRID_FROM_EMAIL")
	}
	if m.fromName == "" {
		m.fromName = os.Getenv("SENDGRID_FROM_NAME")
	}
	if sb := os.Getenv("SENDGRID_SANDBOX"); sb != "" {
		// "true", "1", "on" 등 truthy 처리
		m.sandbox = sb == "true" || sb == "1" || strings.ToLower(sb) == "on"
	}

	if m.apiKey == "" || m.fromEmail == "" {
		return fmt.Errorf("sendgrid config missing: SENDGRID_API_KEY or SENDGRID_FROM_EMAIL")
	}

	from := sgmail.NewEmail(m.fromName, m.fromEmail)
	toAddr := sgmail.NewEmail("", to)

	// 텍스트 바디와 HTML 바디를 모두 포함
	plain := body
	html := "<p>" + htmlEscape(body) + "</p>"

	msg := sgmail.NewSingleEmail(from, subject, toAddr, plain, html)

	// 샌드박스 모드 옵션 (개발/테스트용)
	if m.sandbox {
		msg.MailSettings = &sgmail.MailSettings{
			SandboxMode: &sgmail.Setting{Enable: boolPtr(true)},
		}
	}

	client := sendgrid.NewSendClient(m.apiKey)
	resp, err := client.Send(msg)
	if err != nil {
		return fmt.Errorf("sendgrid send error: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("sendgrid response error: status=%d body=%s", resp.StatusCode, resp.Body)
	}
	return nil
}

// 아주 간단한 HTML 이스케이프
func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		`&`, "&amp;",
		`<`, "&lt;",
		`>`, "&gt;",
		`"`, "&quot;",
		"'", "&#39;",
		"\n", "<br/>",
	)
	return replacer.Replace(s)
}
