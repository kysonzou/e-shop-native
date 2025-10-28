package validator

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"

	"github.com/kyson/e-shop-native/internal/user-srv/biz"
	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"
)

type ValidatorUsecase struct{}

func NewValidator() biz.UserValidator {
	return &ValidatorUsecase{}
}

func (v *ValidatorUsecase) Validate(user *biz.User) error {
	validate, err := getValidator()
	if err != nil {
		return err
	}
	err = validate.Struct(user)
	if err != nil {
		return TranslateValidationError(err)
	}
	return nil
}

var (
	validate *validator.Validate
	once     sync.Once
)

func getValidator() (*validator.Validate, error) {
	var err error
	// validator.New() 是一个昂贵的、一次性的初始化操作。如果在每次HTTP请求或每次需要验证时都去调用 validator.New()，
	// 将会带来巨大且完全不必要的性能开销。
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
	return validate, fmt.Errorf("failed to register validation: %w", err)
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
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)     // 至少一个小写字母
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)     // 至少一个大写字母
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)        // 至少一个数字

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

// 错误判断转换
func TranslateValidationError(err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			return translateFieldError(e)
		}
	}
	return fmt.Errorf("failed to translate validation error: %w", err)
}

func translateFieldError(fe validator.FieldError) error {
	switch fe.Field() {
	case "UserName":
		switch fe.Tag() {
		case "required":
			return apperrors.ErrUsernameFormat.WithMessage("用户名不能为空")
		case "min":
			return apperrors.ErrUsernameFormat.WithMessage("用户名长度不能小于3个字符")
		case "max":
			return apperrors.ErrUsernameFormat.WithMessage("用户名长度不能大于20个字符")
		case "username":
			return apperrors.ErrUsernameFormat.WithMessage("用户名格式错误,支持字母、数字、下划线")
		default:
			return apperrors.ErrUsernameFormat
		}
	case "Password":
		switch fe.Tag() {
		case "required":
			return apperrors.ErrPasswordFormat.WithMessage("密码不能为空")
		case "min":
			return apperrors.ErrPasswordFormat.WithMessage("密码长度不能小于8个字符")
		case "max":
			return apperrors.ErrPasswordFormat.WithMessage("密码长度不能大于64个字符")
		case "password":
			return apperrors.ErrPasswordFormat.WithMessage("密码格式错误,支持字母、数字、特殊字符,且必须包含大小写字母和数字")
		default:
			return apperrors.ErrPasswordFormat
		}
	case "Phone":
		switch fe.Tag() {
		case "required":
			return apperrors.ErrPhoneFormat.WithMessage("手机号不能为空")
		case "phone":
			return apperrors.ErrPhoneFormat
		default:
			return apperrors.ErrPhoneFormat
		}
	case "Email":
		switch fe.Tag() {
		case "required":
			return apperrors.ErrEmailFormat.WithMessage("邮箱不能为空")
		case "email":
			return apperrors.ErrEmailFormat.WithMessage("邮箱格式错误")
		default:
			return apperrors.ErrEmailFormat
		}
	}
	return fmt.Errorf("failed to translate field error: %w", fe)
}
