package batchs

import (
	"context"
	"marketplace-backend/business/batchs"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BatchRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) batchs.Repository {
	return &BatchRepository{
		collection: db.Collection("batchs"),
	}
}

/*
Create
*/

func (br *BatchRepository) Create(domain *batchs.Domain) (batchs.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := br.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return batchs.Domain{}, err
	}

	return *domain, nil
}

/*
Read
*/

func (br *BatchRepository) CountByProposalName(proposalName string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	count, err := br.collection.CountDocuments(ctx, bson.M{
		"name": bson.M{
			"$regex":   proposalName,
			"$options": "i",
		},
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

/*
Update
*/

/*
Delete
*/
