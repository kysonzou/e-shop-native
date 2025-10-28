package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	biz "github.com/kyson/e-shop-native/internal/user-srv/biz"
	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"
)

// 测试用户名格式
func TestUsername(t *testing.T) {
	validate, err := getValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	type TestUser struct {
		UserName string `validate:"required,min=3,max=20,username"`
	}

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"有效用户名-字母", "testuser", false},
		{"有效用户名-数字", "user123", false},
		{"有效用户名-下划线", "test_user", false},
		{"有效用户名-混合", "test_123", false},
		{"无效-太短", "ab", true},
		{"无效-太长", "21_a12345678901234567", true},
		{"无效-包含特殊字符", "test@user", true},
		{"无效-包含空格", "test user", true},
		{"无效-包含中文", "测试用户", true},
		{"无效-空字符串", "", true},
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
		})
	}
}

// 测试密码格式
func TestPassword(t *testing.T) {
	validate, err := getValidator()
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
	}{
		{"有效密码-大小写字母数字", "pAssword123", false},
		{"有效密码-包含特殊字符", "paSs@word1", false},
		{"无效-太短", "shOrtj1", true},
		{"无效-太长", "thisisaverylongpasswordthatexceedsTheMaximumlength1234567890hjygf", true},
		{"无效-仅小写字母", "onlyletters", true},
		{"无效-仅大写字母", "ONLYLETTERS", true},
		{"无效-仅数字", "12345678", true},
		{"无效-小写+大写", "PaSsWopassword", true},
		{"无效-小写+数字", "passw0rd1234", true},
		{"无效-大写+数字", "PASSW0RD1234", true},
		{"无效-包含空格", "pass word1", true},
		{"无效-包含中文", "密码12345", true},
		{"无效-空字符串", "", true},
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
		})
	}
}

// 测试手机号格式
func TestPhone(t *testing.T) {
	validate, err := getValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	type TestUser struct {
		Phone string `validate:"required,phone"`
	}

	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"有效手机号-标准格式", "13800138000", false},
		{"无效-11开头", "11800138000", true},
		{"无效-12开头", "12800138000", true},
		{"无效-太短", "1380013800", true},
		{"无效-太长", "138001380000", true},
		{"无效-包含字母", "13800abc000", true},
		{"无效-包含特殊字符", "13800@#%000", true},
		{"无效-空字符串", "", true},
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
		})
	}
}

// 测试邮箱格式
func TestEmail(t *testing.T) {
	validate, err := getValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	type TestUser struct {
		Email string `validate:"required,email"`
	}

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"有效邮箱-标准格式", "test@example.com", false},
		{"有效邮箱-包含子域", "us23.er!@mail.example.com", false},
		{"有效-包含特殊字符", "test!@example.com", false},
		{"无效-缺少@", "invalidemail.com", true},
		{"无效-缺少用户名", "@example.com", true},
		{"无效-缺少域名", "user@", true},
		{"无效-包含空格", "test @example.com", true},
		{"无效-空字符串", "", true},
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
		})
	}
}

// 测试validator错误转换成业务错误
func TestTranslateValidationError(t *testing.T) {

	type TestUser struct {
		UserName string `validate:"required,min=3,max=20,username"`
		Password string `validate:"required,min=8,max=64,password"`
		Phone    string `validate:"required,phone"`
		Email    string `validate:"required,email"`
	}

	tests := []struct {
		name    string
		user    TestUser
		wantErr error
	}{
		// Add test cases here
		{
			name: "用户名为空",
			user: TestUser{UserName: "",
				Password: "pAssword123",
				Phone:    "13800138000",
				Email:    "test@example.com",
			},
			wantErr: apperrors.ErrUsernameFormat.WithMessage("用户名不能为空"),
		},
		{
			name: "用户名格式错误",
			user: TestUser{UserName: "test@123",
				Password: "pAssword123",
				Phone:    "13800138000",
				Email:    "test@example.com"},
			wantErr: apperrors.ErrUsernameFormat.WithMessage("用户名格式错误,支持字母、数字、下划线"),
		},
		{
			name: "密码为空",
			user: TestUser{UserName: "testuser123",
				Password: "",
				Phone:    "13800138000",
				Email:    "test@example.com"},
			wantErr: apperrors.ErrPasswordFormat.WithMessage("密码不能为空"),
		},
		{
			name: "密码格式错误",
			user: TestUser{UserName: "testuser123",
				Password: "password",
				Phone:    "13800138000",
				Email:    "test@example.com",
			},
			wantErr: apperrors.ErrPasswordFormat.WithMessage("密码格式错误,支持字母、数字、特殊字符,且必须包含大小写字母和数字"),
		},
		{
			name: "手机号为空",
			user: TestUser{UserName: "testuser123",
				Password: "pAssword123",
				Phone:    "",
				Email:    "test@example.com"},
			wantErr: apperrors.ErrPhoneFormat.WithMessage("手机号不能为空"),
		},
		{
			name: "手机号格式错误",
			user: TestUser{UserName: "testuser123",
				Password: "pAssword123",
				Phone:    "1380013800",
				Email:    "test@example.com"},
			wantErr: apperrors.ErrPhoneFormat.WithMessage("手机号格式错误"),
		},
		{
			name: "邮箱为空",
			user: TestUser{UserName: "testuser123",
				Password: "pAssword123",
				Phone:    "13800138000",
				Email:    ""},
			wantErr: apperrors.ErrEmailFormat.WithMessage("邮箱不能为空"),
		},
		{
			name: "邮箱格式错误",
			user: TestUser{UserName: "testuser123",
				Password: "pAssword123",
				Phone:    "13800138000",
				Email:    "test@example"},
			wantErr: apperrors.ErrEmailFormat.WithMessage("邮箱格式错误"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate, err := getValidator()
			if err != nil {
				t.Fatalf("failed to create validator: %v", err)
			}
			err = validate.Struct(tt.user)
			err = TranslateValidationError(err)
			if err.Error() != tt.wantErr.Error() {
				t.Errorf("TranslateValidationError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// GoodPath
func TestValidator(t *testing.T) {
	validate := NewValidator()
	assert.NotNil(t, validate)

	user := &biz.User{
		UserName: "testuser123",
		Password: "pAssword123",
		Phone:    "13800138000",
		Email:    "test@example.com",
	}

	err := validate.Validate(user)
	assert.NoError(t, err)
}
