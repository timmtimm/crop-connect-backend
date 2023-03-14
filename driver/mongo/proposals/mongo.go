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

func (pr *proposalRepository) GetByID(id primitive.ObjectID) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := pr.collection.FindOne(ctx, bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

func (pr *proposalRepository) GetByCommodityID(commodityID primitive.ObjectID) ([]proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := pr.collection.Find(ctx, bson.M{
		"commodityID": commodityID,
		"deletedAt":   bson.M{"$exists": false},
	})
	if err != nil {
		return []proposals.Domain{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []proposals.Domain{}, err
	}

	return ToDomainArray(result), err
}

func (pr *proposalRepository) GetByCommodityIDAndAvailability(commodityID primitive.ObjectID, isAvailable bool) ([]proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := pr.collection.Find(ctx, bson.M{
		"commodityID": commodityID,
		"isAvailable": isAvailable,
		"deletedAt":   bson.M{"$exists": false},
	})
	if err != nil {
		return []proposals.Domain{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []proposals.Domain{}, err
	}

	return ToDomainArray(result), err
}

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

func (pr *proposalRepository) GetByIDAndFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := pr.collection.FindOne(ctx, bson.M{
		"_id":       id,
		"farmerID":  farmerID,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

/*
Update
*/

func (pr *proposalRepository) Update(domain *proposals.Domain) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := pr.collection.UpdateOne(ctx, bson.M{
		"_id":       domain.ID,
		"deletedAt": bson.M{"$exists": false},
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return proposals.Domain{}, err
	}

	return *domain, nil
}

/*
Delete
*/

func (pr *proposalRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := pr.collection.UpdateOne(ctx, bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}, bson.M{
		"$set": bson.M{
			"deletedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
