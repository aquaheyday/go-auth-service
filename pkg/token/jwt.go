package token

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Provider interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateAccessToken(token string) (*JWTClaims, error)
	ValidateRefreshToken(token string) (*JWTClaims, error)
}

// JWTClaims 커스텀 클레임 구조
type JWTClaims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

// 환경 변수에서 시크릿 로드
var (
	accessSecret  = []byte(getEnv("JWT_ACCESS_SECRET", "access-secret-key"))
	refreshSecret = []byte(getEnv("JWT_REFRESH_SECRET", "refresh-secret-key"))
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Access Token: 15분 유효
func GenerateAccessToken(userID string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

// Refresh Token: 7일 유효
func GenerateRefreshToken(userID string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// 토큰 검증 (Access용)
func ValidateAccessToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid access token")
}

// 토큰 검증 (Refresh용)
func ValidateRefreshToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid refresh token")
}

/*import (
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
}*/
