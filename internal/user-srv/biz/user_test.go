package biz_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	biz "github.com/kyson/e-shop-native/internal/user-srv/biz"
	mock "github.com/kyson/e-shop-native/internal/user-srv/biz/mock"
	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"

	gomock "github.com/golang/mock/gomock"
)

// 定义自定义匹配器类型
type userMatcher struct {
	expectedUsername string
	expectedPassword string
	expectedPhone    string
	expectedEmail    string
}

func (m *userMatcher) Matches(x interface{}) bool {
	user, ok := x.(*biz.User)
	if !ok {
		return false
	}
	return user.UserName == m.expectedUsername &&
		user.Password == m.expectedPassword &&
		user.Phone == m.expectedPhone &&
		user.Email == m.expectedEmail
}

func (m *userMatcher) String() string {
	return "user matcher"
}

func NewUserMatcher(username, password, phone, email string) *userMatcher {
	return &userMatcher{
		expectedUsername: username,
		expectedPassword: password,
		expectedPhone:    phone,
		expectedEmail:    email,
	}
}

// 注册用户
func TestUserUsecase_RegisterUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	repo := mock.NewMockUserRepo(ctl)
	validator := mock.NewMockUserValidator(ctl)
	passwordHash := mock.NewMockPasswordHash(ctl)
	uc := biz.NewUserUsecase(repo, validator, passwordHash)

	tests := []struct {
		name      string
		user      *biz.User
		setupMock func(user *biz.User)
		wantErr   error
	}{
		{
			name: "成功注册新用户",
			user: &biz.User{
				ID:       0,
				UserName: "testuser",
				Password: "pAssword123",
				Phone:    "15019458680",
				Email:    "testuser@example.com",
			},
			setupMock: func(user *biz.User) {
				// 验证格式
				validator.EXPECT().Validate(gomock.Eq(user)).Return(nil)
				// 判断是否已经注册
				repo.EXPECT().FindByUsername(gomock.Any(), user.UserName).Return(nil, apperrors.ErrUserNotFound)
				// 密码哈希
				passwordHash.EXPECT().Hash(user.Password).Return("hashed_password", nil)
				// 创建
				repo.EXPECT().Create(gomock.Any(), NewUserMatcher(user.UserName, "hashed_password", user.Phone, user.Email)).DoAndReturn(
					func(ctx context.Context, user *biz.User) (*biz.User, error) {
						user.ID = 1
						return user, nil
					})
			},
			wantErr: nil,
		}, {
			name: "密码格式无效",
			user: &biz.User{
				UserName: "testuser",
				Password: "password",
				Phone:    "15019458680",
				Email:    "testuser@example.com",
			},
			setupMock: func(user *biz.User) {
				// 验证格式
				validator.EXPECT().Validate(gomock.Eq(user)).Return(apperrors.ErrPasswordFormat)
			},
			wantErr: apperrors.ErrPasswordFormat,
		}, {
			name: "用户已存在",
			user: &biz.User{
				UserName: "existinguser",
				Password: "pAssword123",
				Phone:    "15766498680",
				Email:    "existinguser@example.com",
			},
			setupMock: func(user *biz.User) {
				// 验证格式
				validator.EXPECT().Validate(gomock.Eq(user)).Return(nil)
				// 判断是否已经注册
				repo.EXPECT().FindByUsername(gomock.Any(), user.UserName).Return(&biz.User{}, nil)
			},
			wantErr: apperrors.ErrUserAlreadyExists,
		}, {
			name: "密码哈希失败",
			user: &biz.User{
				UserName: "testuser",
				Password: "pAssword123",
				Phone:    "15019458680",
				Email:    "testuser@example.com",
			},
			setupMock: func(user *biz.User) {
				// 验证格式
				validator.EXPECT().Validate(gomock.Eq(user)).Return(nil)
				// 判断是否已经注册
				repo.EXPECT().FindByUsername(gomock.Any(), user.UserName).Return(nil, apperrors.ErrUserNotFound)
				// 密码哈希
				passwordHash.EXPECT().Hash(user.Password).Return("", errors.New("密码哈希失败"))
			},
			wantErr: errors.New("密码哈希失败"),
		}, {
			name: "创建用户失败",
			user: &biz.User{
				UserName: "testuser",
				Password: "pAssword123",
				Phone:    "15019458680",
				Email:    "testuser@example.com",
			},
			setupMock: func(user *biz.User) {
				// 验证格式
				validator.EXPECT().Validate(gomock.Eq(user)).Return(nil)
				// 判断是否已经注册
				repo.EXPECT().FindByUsername(gomock.Any(), user.UserName).Return(nil, apperrors.ErrUserNotFound)
				// 密码哈希
				passwordHash.EXPECT().Hash(user.Password).Return("hashed_password", nil)
				// 创建
				repo.EXPECT().Create(gomock.Any(), NewUserMatcher(user.UserName, "hashed_password", user.Phone, user.Email)).Return(nil, errors.New("创建用户失败"))
			},
			wantErr: errors.New("创建用户失败"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(tt.user)
			_, err := uc.RegisterUser(context.Background(), tt.user)
			assert.Equal(t, err, tt.wantErr)
		})
	}
}

// 登陆
func TestUserUsecase_Login(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	repo := mock.NewMockUserRepo(ctl)
	validate := mock.NewMockUserValidator(ctl)
	passwordHash := mock.NewMockPasswordHash(ctl)
	uc := biz.NewUserUsecase(repo, validate, passwordHash)

	tests := []struct {
		name      string
		username  string
		password  string
		setupMock func(username, password string)
		wantUser  *biz.User
		wantErr   error
	}{
		{
			name:     "成功登录",
			username: "testuser",
			password: "pAssword123",
			setupMock: func(username, password string) {
				repo.EXPECT().FindByUsername(gomock.Any(), username).Return(&biz.User{
					ID: 1, UserName: username, Password: "hashed_password",
				}, nil)
				passwordHash.EXPECT().Virefy(password, "hashed_password").Return(true)
			},
			wantUser: &biz.User{ID: 1, UserName: "testuser", Password: "pAssword123"},
			wantErr:  nil,
		},
		{
			name:     "用户不存在",
			username: "nonexistent",
			password: "password",
			setupMock: func(username, password string) {
				repo.EXPECT().FindByUsername(gomock.Any(), username).Return(nil, apperrors.ErrUserNotFound)
			},
			wantUser: nil,
			wantErr:  apperrors.ErrUserNotFound,
		},
		{
			name:     "密码不正确",
			username: "testuser",
			password: "wrongpassword",
			setupMock: func(username, password string) {
				repo.EXPECT().FindByUsername(gomock.Any(), username).Return(&biz.User{
					ID: 1, UserName: username, Password: "hashed_password",
				}, nil)
				passwordHash.EXPECT().Virefy(password, "hashed_password").Return(false)
			},
			wantUser: nil,
			wantErr:  apperrors.ErrPasswordIncorrect,
		},
		{
			name:     "数据库错误",
			username: "testuser",
			password: "pAssword123",
			setupMock: func(username, password string) {
				repo.EXPECT().FindByUsername(gomock.Any(), username).Return(nil, errors.New("数据库错误"))
			},
			wantUser: nil,
			wantErr:  errors.New("数据库错误"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(tt.username, tt.password)
			user, err := uc.Login(context.Background(), tt.username, tt.password)
			assert.Equal(t, err, tt.wantErr)
			assert.Equal(t, user, tt.wantUser)
		})
	}
}

// 获取用户信息
func TestUserUsecase_GetMyProfile(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	repo := mock.NewMockUserRepo(ctl)
	validate := mock.NewMockUserValidator(ctl)
	passwordHash := mock.NewMockPasswordHash(ctl)
	uc := biz.NewUserUsecase(repo, validate, passwordHash)
	tests := []struct {
		name      string
		userID    uint
		setupMock func(userID uint)
		wantUser  *biz.User
		wantErr   error
	}{
		{
			name:   "成功获取用户信息",
			userID: 1,
			setupMock: func(userID uint) {
				repo.EXPECT().FindByID(gomock.Any(), userID).Return(&biz.User{
					ID: userID, UserName: "testuser", Password: "hashed_pAssword123",
				}, nil)
			},
			wantUser: &biz.User{ID: 1, UserName: "testuser", Password: "hashed_pAssword123"},
			wantErr:  nil,
		}, {
			name:   "用户不存在",
			userID: 2,
			setupMock: func(userID uint) {
				repo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, apperrors.ErrUserNotFound)
			},
			wantUser: nil,
			wantErr:  apperrors.ErrUserNotFound,
		}, {
			name:   "数据库错误",
			userID: 3,
			setupMock: func(userID uint) {
				repo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, errors.New("数据库错误"))
			},
			wantUser: nil,
			wantErr:  errors.New("数据库错误"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(tt.userID)
			user, err := uc.GetMyProfile(context.Background(), tt.userID)
			assert.Equal(t, err, tt.wantErr)
			assert.Equal(t, user, tt.wantUser)
		})
	}
}
