package biz

import (
	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

func NewBcrypt() PasswordHash {
	return &Bcrypt{}
}

func (b *Bcrypt) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", apperrors.ErrPasswordHash
	}
	return string(hashedPassword), nil
}

func (b *Bcrypt) Virefy(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
