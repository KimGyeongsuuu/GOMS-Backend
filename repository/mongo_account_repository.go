package repository

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAccountRepository struct {
	collection *mongo.Collection
}

func NewMongoAccountRepository(db *mongo.Database) *MongoAccountRepository {
	return &MongoAccountRepository{
		collection: db.Collection("accounts"),
	}
}

func (repository *MongoAccountRepository) SaveAccount(ctx context.Context, account *model.Account) error {
	_, err := repository.collection.InsertOne(ctx, account)
	return err
}

func (repository *MongoAccountRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := repository.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *MongoAccountRepository) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account
	err := repository.collection.FindOne(ctx, bson.M{"email": email}).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &account, nil
}

func (repository *MongoAccountRepository) FindByAccountID(ctx context.Context, accountID primitive.ObjectID) (*model.Account, error) {
	var account model.Account
	err := repository.collection.FindOne(ctx, bson.M{"_id": accountID}).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &account, nil
}

func (repository *MongoAccountRepository) FindAllAccount(ctx context.Context) ([]model.Account, error) {
	var accounts []model.Account
	cursor, err := repository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var account model.Account
		if err := cursor.Decode(&account); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (repository *MongoAccountRepository) FindByAccountByStudentInfo(ctx context.Context, searchAccountInput *input.SearchAccountInput) ([]model.Account, error) {
	var accounts []model.Account
	filter := bson.M{}

	if searchAccountInput.Grade != nil {
		filter["grade"] = *searchAccountInput.Grade
	}
	if searchAccountInput.Gender != nil {
		filter["gender"] = *searchAccountInput.Gender
	}
	if searchAccountInput.Name != nil {
		filter["name"] = bson.M{"$regex": *searchAccountInput.Name, "$options": "i"}
	}
	if searchAccountInput.Authority != nil {
		filter["authority"] = *searchAccountInput.Authority
	}
	if searchAccountInput.Major != nil {
		filter["major"] = *searchAccountInput.Major
	}

	cursor, err := repository.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var account model.Account
		if err := cursor.Decode(&account); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (repository *MongoAccountRepository) UpdateAccountAuthority(ctx context.Context, authorityInput *input.UpdateAccountAuthorityInput) error {
	_, err := repository.collection.UpdateOne(ctx, bson.M{"_id": authorityInput.AccountID}, bson.M{"$set": bson.M{"authority": authorityInput.Authority}})
	return err
}

func (repository *MongoAccountRepository) DeleteAccount(ctx context.Context, account *model.Account) error {
	_, err := repository.collection.DeleteOne(ctx, bson.M{"_id": account.ID})
	return err
}
