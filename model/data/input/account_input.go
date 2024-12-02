package input

import (
	"GOMS-BACKEND-GO/model/data/constant"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchAccountInput struct {
	Grade       *int
	Gender      *constant.Gender
	Name        *string
	Authority   *constant.Authority
	IsBlackList *bool
	Major       *constant.Major
}

type UpdateAccountAuthorityInput struct {
	AccountID primitive.ObjectID
	Authority constant.Authority
}
