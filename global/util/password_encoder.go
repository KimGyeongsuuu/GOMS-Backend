package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct{}

func NewPasswordUtil() *Password {
	return &Password{}
}

type UtilPassword interface {
	IsPasswordMatch(rawPassword, encodedPassword string) (bool, error)
	EncodePassword(password string) (string, error)
}

func (p *Password) EncodePassword(password string) (string, error) {
	encodedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("password encode error")
	}
	return string(encodedPassword), nil
}

func (p *Password) IsPasswordMatch(rawPassword string, encodedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(rawPassword))
	if err != nil {
		return false, err
	}
	return true, nil

}
