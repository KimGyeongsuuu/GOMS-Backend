package output

import "GOMS-BACKEND-GO/model/data/constant"

type TokenOutput struct {
	AccessToken     string
	RefreshToken    string
	AccessTokenExp  string
	RefreshTokenExp string
	Authority       constant.Authority
}
