package biz

import (
	"fmt"
	
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

func NewBcrypt() PasswordHash {
	return &Bcrypt{}
}

func (b *Bcrypt) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (b *Bcrypt) Virefy(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
