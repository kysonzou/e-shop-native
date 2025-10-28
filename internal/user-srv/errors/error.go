package errors

import (
	"google.golang.org/grpc/codes"

	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1"
	"github.com/kyson/e-shop-native/pkg/code"
)

// 通用错误
var (
	ErrInternal = code.ErrInternal
)

// 定义用户相关的错误
var (
	ErrUserAlreadyExists = code.New(v1.ErrorCode_USER_ALREADY_EXISTS.String(), "用户已存在", codes.AlreadyExists)
	ErrUserNotFound      = code.New(v1.ErrorCode_USER_NOT_FOUND.String(), "用户不存在", codes.NotFound)

	ErrPasswordIncorrect = code.New(v1.ErrorCode_PASSWORD_INCORRECT.String(), "密码错误", codes.Unauthenticated)
)

// 定义认证相关的错误
var (
	ErrTokenInvalid = code.New(v1.ErrorCode_TOKEN_INVALID.String(), "Token 无效", codes.InvalidArgument)
	ErrTokenExpired = code.New(v1.ErrorCode_TOKEN_EXPIRED.String(), "Token 已过期", codes.Unauthenticated)
)

// 定义验证相关的错误
var (
	ErrUsernameFormat = code.New(v1.ErrorCode_USERNAME_FORMAT_ERROR.String(), "用户名格式错误", codes.InvalidArgument)

	ErrEmailFormat = code.New(v1.ErrorCode_EMAIL_FORMAT_ERROR.String(), "邮箱格式错误", codes.InvalidArgument)

	ErrPhoneFormat = code.New(v1.ErrorCode_PHONE_FORMAT_ERROR.String(), "手机号格式错误", codes.InvalidArgument)

	ErrPasswordFormat = code.New(v1.ErrorCode_PASSWORD_FORMAT_ERROR.String(), "密码格式错误", codes.InvalidArgument)
)
