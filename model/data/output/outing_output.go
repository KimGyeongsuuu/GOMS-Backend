package output

import (
	"GOMS-BACKEND-GO/model/data/constant"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OutingStudentOutput struct {
	AccountID   primitive.ObjectID
	Name        string
	Grade       int
	Major       constant.Major
	Gender      constant.Gender
	ProfileURL  *string
	CreatedTime time.Time
}
