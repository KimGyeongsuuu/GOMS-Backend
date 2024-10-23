package input

import "GOMS-BACKEND-GO/model/data/constant"

type SearchAccountInput struct {
	Grade       *int
	Gender      *constant.Gender
	Name        *string
	Authority   *constant.Authority
	IsBlackList *bool
	Major       *constant.Major
}

type UpdateAccountAuthorityInput struct {
	AccountID uint64
	Authority constant.Authority
}
