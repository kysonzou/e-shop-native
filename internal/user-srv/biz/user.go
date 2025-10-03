package biz

import (
	"context"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint
	UserName string
	Password string
	Phone    string
	Email    string
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrUserNotFound      = errors.New("user not found")
)

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

type userUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) UserService {
	return &userUsecase{
		repo: repo,
	}
}

// RegisterUser registers a new user with the provided details.
func (uc *userUsecase) RegisterUser(ctx context.Context, user *User) (*User, error) {
	// 1. 检查用户名是否已存在
	_, err := uc.repo.FindByUsername(ctx, user.UserName)
	if err == nil {
		return nil, ErrUserAlreadyExists // 用户名已存在
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 其他数据库错误
		return nil, err
	}

	// 2. 校验邮箱格式
	// A simple regex for email validation
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(user.Email) {
		return nil, ErrInvalidEmail
	}

	// 3. 密码hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 2. 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	// 3. 返回用户信息
	return user, nil
}

func (uc *userUsecase) GetMyProfile(ctx context.Context, userID uint) (*User, error) {
	// 1. 获取用户信息
	user, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 2. 返回用户信息
	return user, nil
}
