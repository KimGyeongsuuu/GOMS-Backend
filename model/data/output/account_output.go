package output

import "GOMS-BACKEND-GO/model/data/constant"

type AccountOutput struct {
	AccountID   uint64
	Name        string
	Major       constant.Major
	Grade       int
	Gender      constant.Gender
	ProfileURL  *string
	Authority   constant.Authority
	IsBlackList bool
}

type LateTop3Output struct {
	AccountID  uint64
	Name       string
	Major      constant.Major
	Grade      int
	Gender     constant.Gender
	ProfileURL *string
}
