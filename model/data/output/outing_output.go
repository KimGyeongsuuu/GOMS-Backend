package output

import (
	"GOMS-BACKEND-GO/model/data/constant"
	"time"
)

type OutingStudentOutput struct {
	AccountID   uint64
	Name        string
	Grade       int
	Major       constant.Major
	Gender      constant.Gender
	ProfileURL  *string
	CreatedTime time.Time
}
