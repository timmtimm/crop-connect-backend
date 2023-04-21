package users

import (
	"context"
	"crop_connect/business/users"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) users.Repository {
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

func (ur *UserRepository) GetByQuery(query users.Query) ([]users.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{}

	if query.Name != "" {
		filter["name"] = bson.M{
			"$regex":   query.Name,
			"$options": "i",
		}
	}

	if query.Email != "" {
		filter["email"] = bson.M{
			"$regex":   query.Email,
			"$options": "i",
		}
	}

	if query.PhoneNumber != "" {
		filter["phone_number"] = bson.M{
			"$regex": query.PhoneNumber,
		}
	}

	if query.Role != "" {
		filter["role"] = query.Role
	}

	var result []Model
	cursor, err := ur.collection.Find(ctx, filter)
	if err != nil {
		return []users.Domain{}, 0, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []users.Domain{}, 0, err
	}

	return ToDomainArray(result), len(result), nil
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
