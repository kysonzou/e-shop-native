package data

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/kyson/e-shop-native/internal/user-srv/biz"
	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"
)

type UserPO struct {
	//ID       int64  `gorm:"primaryKey;autoIncrement"`
	UserName string
	Password string
	Email    string
	Phone    string
	gorm.Model
}

func (UserPO) TableName() string {
	return "users"
}

type UserRepo struct {
	data *Data
}

func NewUserRepo(data *Data) biz.UserRepo {
	return &UserRepo{data: data}
}

func (po UserPO) toBizUser() *biz.User {
	return &biz.User{
		ID:       po.ID,
		UserName: po.UserName,
		Password: po.Password,
		Phone:    po.Phone,
		Email:    po.Email,
	}
}

func (r *UserRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	po := &UserPO{
		UserName: user.UserName,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	}
	if err := r.data.db.WithContext(ctx).Create(po).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	user.ID = po.ID
	return user, nil
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	var po UserPO
	if err := r.data.db.WithContext(ctx).Where("user_name = ?", username).First(&po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}
	return po.toBizUser(), nil
}

func (r *UserRepo) FindByID(ctx context.Context, id uint) (*biz.User, error) {
	var po UserPO
	if err := r.data.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return po.toBizUser(), nil
}
