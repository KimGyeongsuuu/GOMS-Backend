package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	return nil, nil
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
		return nil, fmt.Errorf("failed to fetch late students: %w", err)
	}
	defer cursor.Close(ctx)

	var lates []model.Late
	for cursor.Next(ctx) {
		var late model.Late
		if err := cursor.Decode(&late); err != nil {
			return nil, fmt.Errorf("failed to decode late: %w", err)
		}
		lates = append(lates, late)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return lates, nil
}
