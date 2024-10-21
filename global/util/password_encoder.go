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

func IsPasswordMatch(rawPassword string, encodedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(rawPassword))
	if err != nil {
		return false, err
	}
	return true, nil

}
