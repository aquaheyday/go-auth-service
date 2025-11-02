// pkg/token/jwt.go

package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims 커스텀 클레임 구조
type JWTClaims struct {
	UserID  string `json:"sub"`
	TokenID string `json:"jti,omitempty"` // JWT ID 추가
	jwt.RegisteredClaims
}

// 환경 변수에서 시크릿 로드
var (
	accessSecret  = []byte(getEnv("JWT_ACCESS_SECRET", "access-secret-key"))
	refreshSecret = []byte(getEnv("JWT_REFRESH_SECRET", "refresh-secret-key"))
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// 고유 토큰 ID 생성
func generateTokenID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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

// Refresh Token: 7일 유효, 고유 ID 포함
func GenerateRefreshToken(userID string) (string, string, error) {
	tokenID, err := generateTokenID()
	if err != nil {
		return "", "", err
	}

	claims := JWTClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return signedToken, tokenID, nil
}

// Access Token 검증
func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Refresh Token 검증
func ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
