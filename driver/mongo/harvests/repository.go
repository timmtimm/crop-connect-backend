package harvests

import (
	"context"
	"marketplace-backend/business/harvests"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HarvestRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) harvests.Repository {
	return &HarvestRepository{
		collection: db.Collection("harvests"),
	}
}

/*
Create
*/

func (hr *HarvestRepository) Create(domain *harvests.Domain) (harvests.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := hr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return harvests.Domain{}, err
	}

	return *domain, err
}

/*
Read
*/

func (hr *HarvestRepository) GetByBatchID(batchID primitive.ObjectID) (harvests.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result harvests.Domain

	err := hr.collection.FindOne(ctx, bson.M{
		"batchID": batchID,
	}).Decode(&result)
	if err != nil {
		return harvests.Domain{}, err
	}

	return result, nil
}

/*
Update
*/

/*
Delete
*/
