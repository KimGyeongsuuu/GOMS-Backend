package config

import (
	"os"
	"strconv"
)

type JwtProperties struct {
	AccessSecret  []byte
	RefreshSecret []byte
}

type JwtExpTimeProperties struct {
	AccessExp  int64
	RefreshExp int64
}

func LoadJwtProperties() (*JwtProperties, *JwtExpTimeProperties, error) {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	accessExp, err := strconv.ParseInt(os.Getenv("JWT_ACCESS_EXP"), 10, 64)
	if err != nil {
		return nil, nil, err
	}
	refreshExp, err := strconv.ParseInt(os.Getenv("JWT_REFRESH_EXP"), 10, 64)
	if err != nil {
		return nil, nil, err
	}

	return &JwtProperties{
			AccessSecret:  []byte(accessSecret),
			RefreshSecret: []byte(refreshSecret),
		}, &JwtExpTimeProperties{
			AccessExp:  accessExp,
			RefreshExp: refreshExp,
		}, nil
}
