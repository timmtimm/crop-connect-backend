package users

import (
	"context"
	"marketplace-backend/business/users"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) users.Repository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

/*
Create
*/

func (ur *UserRepository) Create(domain *users.Domain) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := ur.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return users.Domain{}, err
	}

	return *domain, err
}

/*
Read
*/

func (ur *UserRepository) GetByID(id primitive.ObjectID) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ur.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return users.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (ur *UserRepository) GetByEmail(email string) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ur.collection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (ur *UserRepository) GetByNameAndRole(name string, role string) ([]users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := ur.collection.Find(ctx, bson.M{
		"name": bson.M{
			"$regex": name,
		},
		"role": role,
	})
	if err != nil {
		return []users.Domain{}, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []users.Domain{}, err
	}

	return ToDomainArray(result), nil
}

/*
Update
*/

func (ur *UserRepository) Update(domain *users.Domain) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := ur.collection.UpdateOne(ctx, bson.M{
		"_id": domain.ID,
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return users.Domain{}, err
	}

	return *domain, nil
}

/*
Delete
*/
