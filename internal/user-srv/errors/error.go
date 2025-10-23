package errors

import (
	"github.com/kyson/e-shop-native/pkg/code"
	"google.golang.org/grpc/codes"
)


// 定义用户相关的错误
var (
	ErrUserAlreadyExists = code.New("USER_ALREADY_EXISTS", "用户已存在", codes.AlreadyExists)
	ErrUserNotFound      = code.New("USER_NOT_FOUND", "用户不存在", codes.NotFound)
	ErrPasswordIncorrect = code.New("PASSWORD_INCORRECT", "密码错误", codes.Unauthenticated)
)

// 定义认证相关的错误
var (
	ErrTokenInvalid = code.New("TOKEN_INVALID", "Token 无效", codes.InvalidArgument)
	ErrTokenExpired = code.New("TOKEN_EXPIRED", "Token 已过期", codes.Unauthenticated)
)

// 定义验证相关的错误
var (
	ErrUsernameFormat = code.New("USERNAME_FORMAT_ERROR", "用户名格式错误", codes.InvalidArgument)

	ErrEmailFormat = code.New("EMAIL_FORMAT_ERROR", "邮箱格式错误", codes.InvalidArgument)

	ErrPhoneFormat = code.New("PHONE_FORMAT_ERROR", "手机号格式错误", codes.InvalidArgument)

	ErrPasswordFormat = code.New("PASSWORD_FORMAT_ERROR", "密码格式错误", codes.InvalidArgument)
)
