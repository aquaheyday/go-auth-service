package token

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Provider interface {
	// GenerateAccessToken은 사용자 식별자를 기반으로 JWT를 만들어 반환합니다.
	GenerateAccessToken(userID string) (string, error)
}

type jwtProvider struct {
	secret []byte
	expiry time.Duration
}

func NewJWTProvider() Provider {
	// JWT_SECRET 와 JWT_EXPIRY(예: "24h") 는 .env 에 설정해 두세요
	exp, err := time.ParseDuration(os.Getenv("JWT_EXPIRY"))
	if err != nil {
		exp = 24 * time.Hour
	}
	return &jwtProvider{
		secret: []byte(os.Getenv("JWT_SECRET")),
		expiry: exp,
	}
}

func (j *jwtProvider) GenerateAccessToken(userID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(j.expiry).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(j.secret)
}
