package model

import "context"

type BlackList struct {
	AccountID uint64
	ExpiredAt int64
}

type BlackListRepository interface {
	SaveBlackList(ctx context.Context, BlackList *BlackList) error
}
