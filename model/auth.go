package model

import (
	"GOMS-BACKEND-GO/model/data/input"
	"context"
)

type AuthUseCase interface {
	SignUp(ctx context.Context, input *input.SignUpInput) error
}
