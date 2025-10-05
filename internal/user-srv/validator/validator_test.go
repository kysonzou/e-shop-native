package validator

import (
	"testing"

	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"

)

func TestUsername(t *testing.T) {
	validate, err := NewValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}


	type TestUser struct {
		UserName string `validate:"required,min=3,max=20,username"`
	}

	tests := []struct {	
		name    string
		username string
		wantErr bool
		err     error
	}{
		{"有效用户名-字母", "testuser", false, nil},
		{"有效用户名-数字", "user123", false, nil},
		{"有效用户名-下划线", "test_user", false, nil},
		{"有效用户名-混合", "test_123", false, nil},
		{"无效-太短", "ab", true, apperrors.ErrUsernameInvalid},
		{"无效-太长", "a123456789012345678901", true, apperrors.ErrUsernameInvalid},
		{"无效-包含特殊字符", "test@user", true, apperrors.ErrUsernameInvalid},
		{"无效-包含空格", "test user", true, apperrors.ErrUsernameInvalid},
		{"无效-包含中文", "测试用户", true, apperrors.ErrUsernameInvalid},
		{"无效-空字符串", "", true, apperrors.ErrUsernameRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &TestUser{
				UserName: tt.username,
			}

			// 判断验证结果是否符合预期
			err := validate.Struct(user)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUsername() error = %v, wantErr = %v", err, tt.wantErr)
			}

			// 判断错误类型是否符合预期
			if (err != nil) != (tt.err != nil) {
				t.Errorf("TranslateValidationError() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil {
				if tansErr := TranslateValidationError(err); tansErr != tt.err {
					t.Errorf("TranslateValidationError() error = %v, wantErr = %v", tansErr, tt.err)
				}
			}
		})
	}
}

func TestPassword(t *testing.T) {
	validate, err := NewValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	type TestUser struct {
		Password string `validate:"required,min=8,max=64,password"`
	}

	tests := []struct {
		name     string
		password string
		wantErr  bool
		err      error
	}{
		{"有效密码-大小写字母数字", "pAssword123", false, nil},
		{"有效密码-包含特殊字符", "paSs@word1", false, nil},
		{"无效-太短", "short1", true, apperrors.ErrPasswordInvalid},
		{"无效-太长", "thisisaverylongpasswordthatexceedsthemaximumlength1234567890", true, apperrors.ErrPasswordInvalid},
		{"无效-仅小写字母", "onlyletters", true, apperrors.ErrPasswordInvalid},
		{"无效-仅大写字母", "ONLYLETTERS", true, apperrors.ErrPasswordInvalid},
		{"无效-仅数字", "12345678", true, apperrors.ErrPasswordInvalid},
		{"无效-小写+大写", "PaSsWopassword", true, apperrors.ErrPasswordInvalid},
		{"无效-小写+数字", "passw0rd1234", true, apperrors.ErrPasswordInvalid},
		{"无效-大写+数字", "PASSW0RD1234", true, apperrors.ErrPasswordInvalid},
		{"无效-包含空格", "pass word1", true, apperrors.ErrPasswordInvalid},
		{"无效-包含中文", "密码12345", true, apperrors.ErrPasswordInvalid},
		{"无效-空字符串", "", true, apperrors.ErrPasswordRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &TestUser{
				Password: tt.password,	
			}
			// 判断验证结果是否符合预期
			err := validate.Struct(user)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestPassword() error = %v, wantErr = %v", err, tt.wantErr)
			}
			// 判断错误类型是否符合预期
			if (err != nil) != (tt.err != nil) {
				t.Errorf("TranslateValidationError() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil {
				if tansErr := TranslateValidationError(err); tansErr != tt.err {
					t.Errorf("TranslateValidationError() error = %v, wantErr = %v", tansErr, tt.err)
				}
			}
		})
	}
}

func TestPhone(t *testing.T) {
	validate, err := NewValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	type TestUser struct {
		Phone    string `validate:"required,phone"`
	}

	tests := []struct {
		name    string
		phone   string
		wantErr bool
		err     error
	}{
		{"有效手机号-标准格式", "13800138000", false, nil},
		{"无效-11开头", "11800138000", true, apperrors.ErrPhoneInvalid},
		{"无效-12开头", "12800138000", true, apperrors.ErrPhoneInvalid},
		{"无效-太短", "1380013800", true, apperrors.ErrPhoneInvalid},
		{"无效-太长", "138001380000", true, apperrors.ErrPhoneInvalid},
		{"无效-包含字母", "13800abc000", true, apperrors.ErrPhoneInvalid},
		{"无效-包含特殊字符", "13800@#%000", true, apperrors.ErrPhoneInvalid},
		{"无效-空字符串", "", true, apperrors.ErrPhoneRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &TestUser{
				Phone: tt.phone,
			}
			// 判断验证结果是否符合预期
			err := validate.Struct(user)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestPhone() error = %v, wantErr = %v", err, tt.wantErr)
			}
			// 判断错误类型是否符合预期
			if (err != nil) != (tt.err != nil) {
				t.Errorf("TranslateValidationError() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil {
				if tansErr := TranslateValidationError(err); tansErr != tt.err {
					t.Errorf("TranslateValidationError() error = %v, wantErr = %v", tansErr, tt.err)
				}
			}
		})
	}
}

func TestEmail(t *testing.T) {
	validate, err := NewValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	type TestUser struct {
		Email    string `validate:"required,email"`
	}

	tests := []struct {
		name    string
		email   string
		wantErr bool
		err     error
	}{
		{"有效邮箱-标准格式", "test@example.com", false, nil},	
		{"有效邮箱-包含子域", "user@mail.example.com", false, nil},
		{"有效-包含特殊字符", "test!@example.com", false, nil},
		{"无效-缺少@", "invalidemail.com", true, apperrors.ErrEmailInvalid},
		{"无效-缺少用户名", "@example.com", true, apperrors.ErrEmailInvalid},
		{"无效-缺少域名", "user@", true, apperrors.ErrEmailInvalid},
		{"无效-包含空格", "test @example.com", true, apperrors.ErrEmailInvalid},
		{"无效-空字符串", "", true, apperrors.ErrEmailRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &TestUser{
				Email: tt.email,
			}
			// 判断验证结果是否符合预期
			err := validate.Struct(user)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestEmail() error = %v, wantErr = %v", err, tt.wantErr)
			}
			// 判断错误类型是否符合预期
			if (err != nil) != (tt.err != nil) {
				t.Errorf("TranslateValidationError() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil {
				if tansErr := TranslateValidationError(err); tansErr != tt.err {
					t.Errorf("TranslateValidationError() error = %v, wantErr = %v", tansErr, tt.err)
				}
			}
		})
	}
}