package output

import (
	"GOMS-BACKEND-GO/model/data/constant"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountOutput struct {
	AccountID   primitive.ObjectID
	Name        string
	Major       constant.Major
	Grade       int
	Gender      constant.Gender
	ProfileURL  *string
	Authority   constant.Authority
	IsBlackList bool
}

type LateOutput struct {
	AccountID  primitive.ObjectID
	Name       string
	Major      constant.Major
	Grade      int
	Gender     constant.Gender
	ProfileURL *string
}
