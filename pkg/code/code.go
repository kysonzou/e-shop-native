package code

import (
	"errors"
	"fmt"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Code interface{
	error
	Code() string
	Message() string
	GrpcCode() codes.Code
	WithMessage(message string) Code
	WithError(err error) Code
}

type ecode struct{
	// 错误码，业务唯一
	code string
	// 错误信息，用于展示给用户，只给出业务错误，内部错不绝对不要暴露给外部
	message string
	// gRPC 错误码，用于 gRPC 协议
	grpccode codes.Code
	// 内部错误
	err error
}

var(
	_errors = make(map[string]Code)
	_mux = sync.Mutex{}
)

func New(code string, message string, grpccode codes.Code) Code {
	_mux.Lock()
	defer _mux.Unlock()

	if _, ok := _errors[code]; ok{
		panic("error code already exists")
	}

	e := &ecode{
		code: code,
		message: message,
		grpccode: grpccode,
	}
	_errors[code] = e
	return e
}

func (e *ecode) Error() string { 
	if e.err != nil {
		return fmt.Sprintf("error: code=%s, message=%s, err=%v", e.code, e.message, e.err)
	}
	return fmt.Sprintf("error: code=%s, message=%s", e.code, e.message)
}
func (e *ecode) Code() string { return e.code }
func (e *ecode) Message() string { return e.message }
func (e *ecode) GrpcCode() codes.Code { return e.grpccode }
func (e *ecode) WithMessage(message string) Code {
	return &ecode{
		code: e.code,
		message: message,
		grpccode: e.grpccode,
		err: e.err,
	}
}
func (e *ecode) WithError(err error) Code {
	return &ecode{
		code: e.code,
		message: e.message,
		grpccode: e.grpccode,
		err: err,
	}
}

// 将一个error 转换为ecode
func FromError(err error) Code {
	// code
	var ec *ecode
	if errors.As(err, &ec) {
		return ec
	}

	// status
	st, ok := status.FromError(err)
	if ok {
		return &ecode{
			code: st.Code().String(),
			message: st.Message(),
			grpccode: st.Code(),
		}
	}

	// 应该设置一个默认的error
	return errInternal
}
// 通用错误
var errInternal = New("INTERNAL_ERROR", "内部服务错误", codes.Internal)