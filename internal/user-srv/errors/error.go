package errors

import "errors"

// 通用错误
var ErrInternal = errors.New("internal error")

// 定义用户相关的错误
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrPasswordHash      = errors.New("password hash error")
)

// 定义认证相关的错误
var (
	ErrJWTInvalid       = errors.New("jwt is invalid")
	ErrJWTParamsInvalid = errors.New("jwt creat failure")
	ErrJWTKeyNotEmpty   = errors.New("jwt key is not empty")
	ErrJWTKeyTooShort   = errors.New("jwt key is too short")
	ErrJWTExpireInvalid = errors.New("jwt expire is invalid")
)

// 定义验证相关的错误
var (
	ErrUsernameRequired = errors.New("username is required")
	ErrUsernameInvalid  = errors.New("username must be 3-20 characters and contain only letters, numbers and underscores")

	ErrEmailRequired = errors.New("email is required")
	ErrEmailInvalid  = errors.New("invalid email format")

	ErrPhoneRequired = errors.New("phone number is required")
	ErrPhoneInvalid  = errors.New("invalid phone number format")

	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordInvalid  = errors.New("password must be 6-20 characters and contain both letters and numbers")
)
