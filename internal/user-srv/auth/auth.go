package auth

import (
	"context"
	"time"

	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kyson/e-shop-native/internal/user-srv/conf"
)

// Claims struct definition
type Claims struct {
	Id       uint   `json:"id"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

type AuthIMP struct {
	jwtKey         []byte
	expireDuration time.Duration
	algorithm      jwt.SigningMethod
	whitelist      []string
}

type Auth interface {
	GenerateToken(ctx context.Context, id uint, userName string) (string, error)
	ParseAndSaveToken(ctx context.Context, tokenS string) (context.Context, error)
	// ToContext(ctx context.Context, claims *Claims) context.Context
	// FromContext(ctx context.Context) (*Claims, bool)
	GetWhiteList() []string
	GetJWTKey() []byte
	GetExpireDuration() time.Duration
	GetAlgorithm() jwt.SigningMethod
}

func NewAuth(c *conf.Auth) (Auth, error) {
	if c.JwtKey == "" {
		return nil, apperrors.ErrJWTKeyNotEmpty
	}
	if len(c.JwtKey) < 32 {
		return nil, apperrors.ErrJWTKeyTooShort
	}
	if c.ExpireDuration <= 0 {
		return nil, apperrors.ErrJWTExpireInvalid
	}
	var algorithm jwt.SigningMethod
	switch c.Algorithm {
	case "HS256":
		algorithm = jwt.SigningMethodHS256
	case "HS384":
		algorithm = jwt.SigningMethodHS384
	case "HS512":
		algorithm = jwt.SigningMethodHS512
	default:
		algorithm = jwt.SigningMethodHS256
	}
	authIMP := &AuthIMP{
		jwtKey:         []byte(c.JwtKey),
		expireDuration: time.Second * time.Duration(c.ExpireDuration),
		algorithm:      algorithm,
		whitelist:      c.Whitelist,
	}
	return authIMP, nil
}

func (auth *AuthIMP) GenerateToken(ctx context.Context, id uint, userName string) (string, error) {
	if id <= 0 || userName == "" {
		return "", apperrors.ErrJWTParamsInvalid
	}
	expirationTime := time.Now().Add(auth.expireDuration) // 每次生成时计算
	claims := Claims{
		Id:       id,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(auth.algorithm, claims)  // 计算claims "指纹"
	tokenString, err := token.SignedString(auth.jwtKey) // 密钥签名
	return tokenString, err
}

// ParseToken 解析并验证一个 JWT 字符串
func (auth *AuthIMP) ParseAndSaveToken(ctx context.Context, tokenS string) (context.Context, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenS, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrJWTInvalid
		}
		return auth.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return ctx, apperrors.ErrJWTInvalid
	}

	return ToContext(ctx, claims), nil
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

func (a *AuthIMP) GetWhiteList() []string {
	return a.whitelist
}

func (a *AuthIMP) GetJWTKey() []byte {
	return a.jwtKey
}

func (a *AuthIMP) GetExpireDuration() time.Duration {
	return a.expireDuration
}

func (a *AuthIMP) GetAlgorithm() jwt.SigningMethod {
	return a.algorithm
}
