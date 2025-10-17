package biz_test

import (
	"testing"

	biz "github.com/kyson/e-shop-native/internal/user-srv/biz"
)

func TestBcrypt(t *testing.T) {
	bcrypt := biz.NewBcrypt()
	hashedPassword, err := bcrypt.Hash("password")
	if err != nil {
		t.Fatal(err)
	}
	ok := bcrypt.Virefy("password", hashedPassword)
	if !ok {
		t.Fatal("Virefy failed")
	}
}
