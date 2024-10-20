package util

import (
	"golang.org/x/crypto/bcrypt"
)

func EncodePassword(password string) (string, error) {
	encodedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encodedPassword), nil
}
