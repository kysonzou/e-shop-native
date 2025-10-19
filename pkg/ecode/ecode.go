package ecode

import (
	"errors"
	"fmt"

	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1"
	"google.golang.org/grpc/codes"
)

type ECode interface{
	error
	Code() int
	Message() string
	GrpcCode() codes.Code
	Detail() *v1.UserErr
}

type bussinessError struct{
	code int 
	message string
	grpccode codes.Code
	detail *v1.UserErr
}

func New(code int, message string, grpccode codes.Code) ECode {
	return &bussinessError{
		code: code,
		message: message,
		grpccode: grpccode,
		detail: &v1.UserErr{
			Code: int32(code),
			Message: message,
		},
	}
}

func (e *bussinessError) Code() int { return e.code }
func (e *bussinessError) Message() string { return e.message}
func (e *bussinessError) GrpcCode() codes.Code { return e.grpccode}
func (e *bussinessError) Detail() *v1.UserErr { return e.detail}
func (e *bussinessError) Error() string { 
	return fmt.Sprintf("ecode: code=%d, message=%s", e.code, e.message)
}


// 将一个error 转换为ecode
func FromError(err error) (ECode, bool) {
	// 使用 errors.As 进行接口断言
	var ec ECode
	if errors.As(err, &ec) {
		return ec, true
	}
	// 如果不是 ECode，返回一个默认的内部错误
	return New(-1, "internal server error", codes.Internal), false
}