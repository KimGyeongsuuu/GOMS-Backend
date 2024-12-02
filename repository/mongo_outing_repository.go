package repository

import (
	"GOMS-BACKEND-GO/global/error/status"
	"GOMS-BACKEND-GO/model"
	"context"
	"net/http"

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
		return nil, status.NewError(http.StatusNotFound, "no outings student")
	}

	return outings, nil
}

func (repository *MongoOutingRepository) FindByOutingAccountNameContaining(ctx context.Context, name string) ([]model.Outing, error) {
	var outings []model.Outing

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.M{
			"from":         "accounts",   // accounts 컬렉션과 조회를 한다
			"localField":   "account_id", // account_id 필드 값과 accounts 컬렉션의 _id가 같은 데이터를 조회 한다.
			"foreignField": "_id",
			"as":           "account", // 위의 조건에 맞는 데이터를 account라는 필드에 저장 (이름 직접 지정)
		}}},
		{{"$unwind", bson.M{"path": "$account"}}},                    // 배열 형태로 account에 정보를 담는다.
		{{"$match", bson.M{"account.name": bson.M{"$regex": name}}}}, // account의 name과 입력받은 name을 비교해서 데이터 조회
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
		return nil, status.NewError(http.StatusNotFound, "no outings found for this account name")
	}

	return outings, nil
}
