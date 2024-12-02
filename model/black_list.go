package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlackList struct {
	AccountID primitive.ObjectID
	ExpiredAt int64
}

type BlackListRepository interface {
	SaveBlackList(ctx context.Context, blackList *BlackList) error
	DeleteBlackList(ctx context.Context, blackList *BlackList) error
	FindBlackListByAccountID(ctx context.Context, accountID primitive.ObjectID) (*BlackList, error)
	ExistsByAccountID(ctx context.Context, accountID primitive.ObjectID) (bool, error)
}
