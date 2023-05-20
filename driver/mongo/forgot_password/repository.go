package forgot_password

import (
	"context"
	forgotPassword "crop_connect/business/forgot_password"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ForgotPasswordRepository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) forgotPassword.Repository {
	return &ForgotPasswordRepository{
		collection: db.Collection("forgotPasswords"),
	}
}

/*
Create
*/

func (fpr *ForgotPasswordRepository) Create(domain *forgotPassword.Domain) (forgotPassword.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := fpr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return forgotPassword.Domain{}, err
	}

	return *domain, err
}

/*
Read
*/

func (fr *ForgotPasswordRepository) GetByToken(token string) (forgotPassword.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := fr.collection.FindOne(ctx, bson.M{
		"token": token,
	}).Decode(&result)
	if err != nil {
		return forgotPassword.Domain{}, err
	}

	return result.ToDomain(), nil
}

/*
Update
*/

func (fr *ForgotPasswordRepository) Update(domain *forgotPassword.Domain) (forgotPassword.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := fr.collection.UpdateOne(ctx, bson.M{
		"_id": domain.ID,
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return forgotPassword.Domain{}, err
	}

	return *domain, nil
}

/*
Delete
*/

func (fpr *ForgotPasswordRepository) HardDelete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := fpr.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
