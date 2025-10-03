package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
var ErrTokenInvalid = errors.New("token is invalid")

var jwtKey = []byte("ahkPzSJ6auFD2WZHt5NFfixFSI3JmXm4isbTs8y29Zs=")
var expirationTime = time.Now().Add(24 * time.Hour)

// Claims struct definition
type Claims struct {
	Id              uint   `json:"id"`
	UserName        string `json:"userName"`
	jwt.RegisteredClaims
}

func GenerateToken(id uint, userName string) (string, error) {
	claims := Claims{
		Id:       id,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 计算claims "指纹"
	tokenString, err := token.SignedString(jwtKey) // 密钥签名
	return tokenString, err
}

// ParseToken 解析并验证一个 JWT 字符串
func ParseToken(tokenS string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenS, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

type claimKey struct{}

// ToContext 和 FromContext
func ToContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimKey{}, claims)
}

func FromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimKey{}).(*Claims)
	return claims, ok
}