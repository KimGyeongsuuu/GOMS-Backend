package repository

import (
	"GOMS-BACKEND-GO/global/error/status"
	"GOMS-BACKEND-GO/model"
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoLateRepository struct {
	collection *mongo.Collection
}

func NewMongoLateRepository(db *mongo.Database) *MongoLateRepository {
	return &MongoLateRepository{
		collection: db.Collection("lates"),
	}
}

func (repository *MongoLateRepository) FindTop3ByOrderByAccountDesc(ctx context.Context) ([]model.Late, error) {
	var lates []model.Late

	pipeline := mongo.Pipeline{
		{{"$group", bson.M{
			"_id":   "$account_id",              // account_id로 그룹화
			"count": bson.M{"$sum": 1},          // 각 account_id의 출현 횟수
			"late":  bson.M{"$first": "$$ROOT"}, // 해당 account_id의 첫 번째 Late 정보를 가져옵니다.
		}}},
		{{"$sort", bson.M{"count": -1}}}, // 출현 횟수(count) 기준으로 내림차순 정렬
		{{"$limit", 3}},                  // 상위 3개 결과만 선택
	}

	cursor, err := repository.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch top lates: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			AccountID primitive.ObjectID `bson:"_id"`
			Late      model.Late         `bson:"late"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		lates = append(lates, result.Late)
	}

	return lates, nil
}

func (repository *MongoLateRepository) FindLateByCreatedAt(ctx context.Context, date time.Time) ([]model.Late, error) {
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	cursor, err := repository.collection.Find(ctx, filter)
	if err != nil {
		return nil, status.NewError(http.StatusBadRequest, "failed to fetch late students")
	}
	defer cursor.Close(ctx)

	var lates []model.Late
	for cursor.Next(ctx) {
		var late model.Late
		if err := cursor.Decode(&late); err != nil {
			return nil, status.NewError(http.StatusBadRequest, "failed to decode late")
		}
		lates = append(lates, late)
	}

	if err := cursor.Err(); err != nil {
		return nil, status.NewError(http.StatusInternalServerError, "cursor error")
	}

	return lates, nil
}
