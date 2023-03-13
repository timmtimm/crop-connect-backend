package proposals

import (
	"context"
	"marketplace-backend/business/proposals"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type proposalRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) proposals.Repository {
	return &proposalRepository{
		collection: db.Collection("proposals"),
	}
}

/*
Create
*/

func (pr *proposalRepository) Create(domain *proposals.Domain) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := pr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return proposals.Domain{}, err
	}

	return *domain, err
}

/*
Read
*/

func (pr *proposalRepository) GetByCommodityIDAndName(commodityID primitive.ObjectID, name string) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := pr.collection.FindOne(ctx, bson.M{
		"commodityID": commodityID,
		"name":        name,
		"deletedAt":   bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

/*
Update
*/

/*
Delete
*/
