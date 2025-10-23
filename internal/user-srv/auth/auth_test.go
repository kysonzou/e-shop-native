package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	//"github.com/kyson/e-shop-native/internal/user-srv/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"
)

// type fakeAuthOptions func (*conf.Auth);

// func withFakeAuthJWTKey(jwtKey string) fakeAuthOptions {
// 	return func (a *conf.Auth) {
// 		a.JwtKey = jwtKey
// 	}
// }
// func withFakeAuthAlgorithm(algorithm string) fakeAuthOptions {
// 	return func (a *conf.Auth) () {
// 		a.Algorithm = algorithm
// 	}
// }
// func withFakeAuthExpireDuration(expireDuration int64) fakeAuthOptions {
// 	return func (a *conf.Auth) () {
// 		a.ExpireDuration = expireDuration
// 	}
// }
// func withFakeAuthWhitelist(whitelist []string) fakeAuthOptions {
// 	return func (a *conf.Auth) () {
// 		a.Whitelist = whitelist
// 	}
// }

// func getFakeAuthConfig(opts...fakeAuthOptions) *conf.Auth{
// 	authConfig := &conf.Auth{
// 		Algorithm: "HS256",
// 		JwtKey: "ahkPzSJ6auFD2WZHt5NFfixFSI3JmXm4isbTs8y29Zs=",
// 		ExpireDuration: 3600,
// 		Whitelist: []string{"/user/register", "/user/login"},
// 	}
// 	for _, opt := range opts {
// 		opt(authConfig)
// 	}
// 	return authConfig
// }

// // 测试confing的动态配置
// func TestFakeAuthConfig(t *testing.T){
// 	jwt_key := "123"
// 	expireDuration := int64(25)
// 	algorithm := "Hash512"
// 	whiteList := []string{"/user/register"}
// 	conf := getFakeAuthConfig(withFakeAuthJWTKey(jwt_key),
// 							  withFakeAuthExpireDuration(expireDuration),
// 							  withFakeAuthAlgorithm(algorithm),
// 							  withFakeAuthWhitelist(whiteList))
// 	assert.Equal(t, conf.JwtKey, jwt_key)
// 	assert.Equal(t, conf.ExpireDuration, expireDuration)
// 	assert.Equal(t, conf.Algorithm, algorithm)
// 	assert.Equal(t, conf.Whitelist, whiteList)
// }

// 测试NewAuth函数
func TestGetter(t *testing.T) {
	tests := []struct {
		name           string
		config         *conf.Auth
		Algorithm      jwt.SigningMethod
		ExpireDuration time.Duration
		wantAuth       bool
		wantErr        error
	}{
		{
			name: "正常生成auth",
			config: &conf.Auth{
				Algorithm:      "HS512",
				JwtKey:         "ahkPzSJ6auFD2WZHt5NFfixFSI3JmXm4isbTs8y29Zs=",
				ExpireDuration: 3600,
				Whitelist:      []string{"/user/register", "/user/login"},
			},
			Algorithm:      jwt.SigningMethodHS512,
			ExpireDuration: time.Second * time.Duration(3600),
		},
	}

	for _, tt := range tests {
		auth := NewAuth(tt.config)
		assert.Equal(t, tt.config.JwtKey, string(auth.GetJWTKey()))
		assert.Equal(t, tt.Algorithm, auth.GetAlgorithm())
		assert.Equal(t, tt.ExpireDuration, auth.GetExpireDuration())
		assert.Equal(t, tt.config.Whitelist, auth.GetWhiteList())
		
	}
}

// 测试生成Token函数
func TestGenerateToken(t *testing.T) {
	config := &conf.Auth{
		Algorithm:      "HS256",
		JwtKey:         "ahkPzSJ6auFD2WZHt5NFfixFSI3JmXm4isbTs8y29Zs=",
		ExpireDuration: 3600,
		Whitelist:      []string{"/user/register", "/user/login"},
	}

	tests := []struct {
		name      string
		uid       uint
		userName  string
		config    *conf.Auth
		wantToken bool
		wantErr   error
	}{
		{
			name:      "正常生成token",
			uid:       1,
			userName:  "testUser",
			config:    config,
			wantToken: true,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		auth := NewAuth(tt.config)
		token, err := auth.GenerateToken(context.Background(), tt.uid, tt.userName)
		assert.Equal(t, tt.wantToken, token != "")
		assert.Equal(t, tt.wantErr, err)
	}
}

// 测试Token解析和获取
func TestParseToken(t *testing.T) {
	config := &conf.Auth{
		Algorithm:      "HS512",
		JwtKey:         string("ahkPzSJ6auFD2WZHt5NFfixFSI3JmXm4isbTs8y29Zs="),
		ExpireDuration: 600,
		Whitelist:      []string{"/user/register", "/user/login"},
	}

	// 正常的tokenString
	validClaims := Claims{
		Id:       1,
		UserName: "testUser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(config.ExpireDuration))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims)
	tokenString, err := token.SignedString([]byte(config.JwtKey))
	assert.NoError(t, err)

	// 修改过的tokenString
	editTokenString := tokenString + "1"

	// 过期的Token
	expireClaims := validClaims
	expireClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(-100)))
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, expireClaims)
	expireTokenString, err := token.SignedString([]byte(config.JwtKey))
	assert.NoError(t, err)

	tests := []struct {
		name        string
		tokenString string
		wantErr     error
	}{
		{
			name:        "正常解析",
			tokenString: tokenString,
			wantErr:     nil,
		}, {
			name:        "被修改的token解析失败",
			tokenString: editTokenString,
			wantErr:     apperrors.ErrTokenInvalid,
		}, {
			name:        "过期的token解析失败",
			tokenString: expireTokenString,
			wantErr:     apperrors.ErrTokenInvalid,
		},
	}

	for _, tt := range tests {
		auth := NewAuth(config)

		ctx, err := auth.ParseAndSaveToken(context.Background(), tt.tokenString)
		assert.Equal(t, tt.wantErr, err)
		if err == nil {
			claims, ok := FromContext(ctx) // 验证读取信息
			assert.True(t, ok)
			assert.Equal(t, validClaims.Id, claims.Id)
			assert.Equal(t, validClaims.UserName, claims.UserName)
			assert.Equal(t, validClaims.ExpiresAt, claims.ExpiresAt)
		}
	}
}

// GoodPath
func TestPath(t *testing.T) {
	config := conf.Auth{
		Algorithm:      "HS512",
		JwtKey:         string("ahkPzSJ6auFD2WZHt5NFfixFSI3JmXm4isbTs8y29Zs="),
		ExpireDuration: 3600,
		Whitelist:      []string{"/user/register", "/user/login"},
	}
	auth := NewAuth(&config)

	tokenString, err := auth.GenerateToken(context.Background(), 1, "testUser")
	assert.NoError(t, err)

	ctx, err := auth.ParseAndSaveToken(context.Background(), tokenString)
	assert.NoError(t, err)

	claims, ok := FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, uint(1), claims.Id)
	assert.Equal(t, "testUser", claims.UserName)
	expectedExpiresAt := time.Now().Add(time.Second * time.Duration(config.ExpireDuration))
	assert.WithinDuration(t, expectedExpiresAt, claims.ExpiresAt.Time, 2*time.Second) // 允许2秒的误差
}
