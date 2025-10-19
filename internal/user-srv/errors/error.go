package errors

import (
	"github.com/kyson/e-shop-native/pkg/ecode"
	"google.golang.org/grpc/codes"
)

// 假设我们定义一个7位错误码 A-BB-CCC：
// A (1位): 错误级别
// 1: 客户端错误 (Client-side)
// 2: 服务端错误 (Server-side)
// 3: 第三方依赖错误 (Third-party)
// BB (2位): 服务/模块ID
// 01: 用户服务
// 02: 商品服务
// 03: 订单服务
// CCC (3位): 具体错误序列号，在该模块内自增。


// 通用错误
var ErrInternal = ecode.New(202001, "internal server error", codes.Internal) 

// 定义用户相关的错误
var (
	ErrUserAlreadyExists = ecode.New(202002, "user already exists", codes.AlreadyExists)  
	ErrUserNotFound      = ecode.New(202003, "user not found", codes.NotFound)   
	ErrPasswordHash      = ecode.New(202004, "password hash error", codes.InvalidArgument)   
)

// 定义认证相关的错误
var (
	ErrJWTInvalid       = ecode.New(202005, "invalid token", codes.Unauthenticated)
	ErrJWTParamsInvalid = ecode.New(202006, "invalid params", codes.InvalidArgument)
	ErrJWTKeyNotEmpty   = ecode.New(202007, "jwt key is empty", codes.InvalidArgument)
	ErrJWTKeyTooShort   = ecode.New(202008, "jwt key is too short", codes.InvalidArgument)
	ErrJWTExpireInvalid = ecode.New(202009, "jwt expire is invalid", codes.InvalidArgument)
)

// 定义验证相关的错误
var (
	ErrUsernameRequired = ecode.New(202010, "username is required", codes.InvalidArgument)
	ErrUsernameInvalid  = ecode.New(202011, "invalid username format", codes.InvalidArgument)

	ErrEmailRequired = ecode.New(202012, "email is required", codes.InvalidArgument)
	ErrEmailInvalid  = ecode.New(202013, "invalid email format", codes.InvalidArgument)

	ErrPhoneRequired = ecode.New(202014, "phone number is required", codes.InvalidArgument)
	ErrPhoneInvalid  = ecode.New(202015, "invalid phone number format", codes.InvalidArgument)

	ErrPasswordRequired = ecode.New(202016, "password is required", codes.InvalidArgument)	
	ErrPasswordInvalid  = ecode.New(202017, "invalid password format", codes.InvalidArgument)
)

