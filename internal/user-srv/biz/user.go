package biz

//go:generate mockgen -source=user.go -destination=mock/mocker_user.go -package=mock

import (
	"context"
	"errors"
	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"
)

type User struct {
	ID       uint
	UserName string `validate:"required,min=3,max=20,username"`
	Password string `validate:"required,min=8,max=64,password"`
	Phone    string `validate:"required,phone"`
	Email    string `validate:"required,email"`
}

type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
}

type UserService interface {
	RegisterUser(ctx context.Context, user *User) (*User, error)
	Login(ctx context.Context, username, password string) (*User, error)
	GetMyProfile(ctx context.Context, userID uint) (*User, error)
}

// 验证用户信息是否符合要求
type UserValidator interface {
	Validate(user *User) error
}

type PasswordHash interface {
	Hash(password string) (string, error)
	Virefy(password, hashedPassword string) bool
}

type userUsecase struct {
	repo      UserRepo
	validator UserValidator
	bcrypt    PasswordHash
}

func NewUserUsecase(repo UserRepo, validator UserValidator, bcrypt PasswordHash) UserService {
	return &userUsecase{
		repo:      repo,
		validator: validator,
		bcrypt:    bcrypt,
	}
}

// RegisterUser registers a new user with the provided details.
func (uc *userUsecase) RegisterUser(ctx context.Context, user *User) (*User, error) {
	// 1. 校验格式（用户名、邮箱、密码、手机号）
	if err := uc.validator.Validate(user); err != nil {
		return nil, err
	}

	// 2. 检查用户名是否已存在
	_, err := uc.repo.FindByUsername(ctx, user.UserName)
	if err == nil {
		return nil, apperrors.ErrUserAlreadyExists // 用户名已存在
	}
	if !errors.Is(err, apperrors.ErrUserNotFound) { // 非用户不存在错误
		// 其他数据库错误
		return nil, err
	}

	// 3. 密码hash
	ps, err := uc.bcrypt.Hash(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = ps

	// 4. 创建新用户
	createdUser, err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

// Login authenticates a user with the provided username and password.
func (uc *userUsecase) Login(ctx context.Context, username, password string) (*User, error) {
	// 1. 获取用户信息
	user, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	// 2. 验证密码
	ok := uc.bcrypt.Virefy(password, user.Password)
	if !ok {
		return nil, apperrors.ErrPasswordInvalid
	}

	user.Password = password
	// 3. 返回用户信息
	return user, nil
}

func (uc *userUsecase) GetMyProfile(ctx context.Context, userID uint) (*User, error) {
	// 1. 获取用户信息
	user, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	// 2. 返回用户信息
	return user, nil
}
