package validator

import (
	"regexp"
	"sync"
	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"
	"github.com/go-playground/validator/v10"
)	

var(
	validate *validator.Validate
	once sync.Once
)

func NewValidator() (*validator.Validate, error) {
	var err error	
	once.Do(func() {
		validate = validator.New()
		// 注册自定义验证器
		err = validate.RegisterValidation("username", ValidateUsername)
		if err != nil {
			return
		}
		err = validate.RegisterValidation("password", ValidatePassword)
		if err != nil {
			return
		}
		err = validate.RegisterValidation("phone", ValidatePhone)
	})
	return validate, err
}

func ValidateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	
	// 只允许字母、数字和下划线，长度3-20个字符
	// if len(username) < 3 || len(username) > 20 {
	// 	return false
	// }
	//regular := `^[a-zA-Z0-9_]{3,20}$` // 只允许字母、数字和下划线，限制大小
	regular := `^[a-zA-Z0-9_]+$` // 只允许字母、数字和下划线
	matched, err := regexp.MatchString(regular, username)
	if err != nil {
		return false
	}
	return matched
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	// 密码至少8个字符，包含至少一个大写字母、一个小写字母和一个数字
	// if len(password) < 8 || len(password) > 64 {
	// 	return false
	// }

	//regular := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,64}$` 
	
	// Go's regexp does not support lookahead, so check conditions manually
	validChars := regexp.MustCompile(`^[\S]+$`).MatchString(password) // 只允许字母、数字、特殊字符，且不包含空白字符
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password) // 至少一个小写字母
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password) // 至少一个大写字母
	hasDigit := regexp.MustCompile(`\d`).MatchString(password) // 至少一个数字

	if !hasLower || !hasUpper || !hasDigit || !validChars {
		return false
	}
	return true
}

func ValidatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	
	// 简单的手机号格式验证（中国手机号）
	regular := `^1[3-9]\d{9}$`
	matched, err := regexp.MatchString(regular, phone)
	if err != nil {
		return false
	}
	return matched
}

// 错误判断
func TranslateValidationError(err error) error {
	if err == nil {
		return nil
	}
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			return translateFieldError(e)
		}
	}
	return err
}

func translateFieldError(fe validator.FieldError) error {
	switch fe.Field() {
	case "UserName":
		switch fe.Tag() {
		case "required":
			return apperrors.ErrUsernameRequired
		case "username":
			return apperrors.ErrUsernameInvalid
		default:
			return apperrors.ErrUsernameInvalid
		}
	case "Password":
		switch fe.Tag() {
		case "required":
			return apperrors.ErrPasswordRequired
		case "password":
			return apperrors.ErrPasswordInvalid
		default:
			return apperrors.ErrPasswordInvalid
		}
	case "Phone":
		switch fe.Tag() {
		case "required":
			return apperrors.ErrPhoneRequired
		case "len":
			return apperrors.ErrPhoneInvalid
		case "phone":
			return apperrors.ErrPhoneInvalid
		default:
			return apperrors.ErrPhoneInvalid
		}
	case "Email":
		switch fe.Tag() {
		case "required":
			return apperrors.ErrEmailRequired
		case "email":
			return apperrors.ErrEmailInvalid
		default:
			return apperrors.ErrEmailInvalid
		}	
	}	
	return fe
}