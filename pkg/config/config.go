package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	GRPCPort          string
	DatabaseURL       string
	RedisAddr         string
	SMTPHost          string
	SMTPPort          int
	SMTPUser          string
	SMTPPass          string
	LogLevel          string
	SendGridAPIKey    string
	SendGridFromEmail string
	SendGridFromName  string
	SendGridSandbox   bool
}

func LoadConfig() (*Config, error) {
	// 1) .env 파일 로드 (실패해도 무시)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, relying on ENV vars: %v", err)
	}

	// 2) viper 설정
	viper.AutomaticEnv()        // ENV 변수 우선
	viper.SetConfigFile(".env") // .env 파일도 읽기
	_ = viper.ReadInConfig()    // 파일 없으면 에러 무시

	cfg := &Config{
		GRPCPort:          viper.GetString("GRPC_PORT"),
		DatabaseURL:       viper.GetString("DATABASE_URL"),
		RedisAddr:         viper.GetString("REDIS_ADDR"),
		SMTPHost:          viper.GetString("SMTP_HOST"),
		SMTPPort:          viper.GetInt("SMTP_PORT"),
		SMTPUser:          viper.GetString("SMTP_USER"),
		SMTPPass:          viper.GetString("SMTP_PASS"),
		LogLevel:          viper.GetString("LOG_LEVEL"),
		SendGridAPIKey:    viper.GetString("SENDGRID_API_KEY"),
		SendGridFromEmail: viper.GetString("SENDGRID_FROM_EMAIL"),
		SendGridFromName:  viper.GetString("SENDGRID_FROM_NAME"),
		SendGridSandbox:   viper.GetBool("SENDGRID_SANDBOX"),
	}

	return cfg, nil
}
