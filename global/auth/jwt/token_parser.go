package jwt

import (
	"errors"
	"strings"
)

type TokenParser struct{}

func NewTokenParser() *TokenParser {
	return &TokenParser{}
}

func (jp *TokenParser) ParseRefreshToken(refreshToken string) (string, error) {
	if strings.HasPrefix(refreshToken, "Bearer ") {
		return strings.TrimPrefix(refreshToken, "Bearer "), nil
	}
	return "", errors.New("invalid refresh token format")
}
