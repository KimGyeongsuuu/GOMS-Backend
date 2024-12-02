package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoOutingRepository struct {
	collection *mongo.Collection
}

func NewMongoOutingRepository(db *mongo.Database) *MongoOutingRepository {
	return &MongoOutingRepository{
		collection: db.Collection("outings"),
	}
}

func (repository *MongoOutingRepository) SaveOutingStudent(ctx context.Context, outing *model.Outing) error {
	_, err := repository.collection.InsertOne(ctx, outing)
	return err
}

func (repository *MongoOutingRepository) ExistsOutingByAccountID(ctx context.Context, accountID primitive.ObjectID) (bool, error) {
	count, err := repository.collection.CountDocuments(ctx, bson.M{"account_id": accountID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *MongoOutingRepository) DeleteOutingByAccountID(ctx context.Context, accountID primitive.ObjectID) error {
	_, err := repository.collection.DeleteOne(ctx, bson.M{"account_id": accountID})
	return err
}

func (repository *MongoOutingRepository) FindAllOuting(ctx context.Context) ([]model.Outing, error) {
	var outings []model.Outing
	cursor, err := repository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var outing model.Outing
		if err := cursor.Decode(&outing); err != nil {
			return nil, err
		}
		outings = append(outings, outing)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(outings) == 0 {
		return nil, errors.New("no outings found")
	}

	return outings, nil
}

func (repository *MongoOutingRepository) FindByOutingAccountNameContaining(ctx context.Context, name string) ([]model.Outing, error) {
	var outings []model.Outing

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.M{
			"from":         "accounts",
			"localField":   "account_id",
			"foreignField": "_id",
			"as":           "account",
		}}},
		{{"$unwind", bson.M{"path": "$account"}}},
		{{"$match", bson.M{"account.name": bson.M{"$regex": name}}}},
	}

	cursor, err := repository.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var outing model.Outing
		if err := cursor.Decode(&outing); err != nil {
			return nil, err
		}
		outings = append(outings, outing)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(outings) == 0 {
		return nil, errors.New("no outings found for this account name")
	}

	return outings, nil
}
