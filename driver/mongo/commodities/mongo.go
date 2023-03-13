package commodities

import (
	"context"
	"fmt"
	"marketplace-backend/business/commodities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type commoditiesRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) commodities.Repository {
	return &commoditiesRepository{
		collection: db.Collection("commodities"),
	}
}

/*
Create
*/

func (cr *commoditiesRepository) Create(domain *commodities.Domain) (commodities.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := cr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return commodities.Domain{}, err
	}

	return *domain, err
}

/*
Read
*/

func (cr *commoditiesRepository) GetByIDAndFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (commodities.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := cr.collection.FindOne(ctx, bson.M{
		"_id":       id,
		"farmerID":  farmerID,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

func (cr *commoditiesRepository) GetByName(name string) (commodities.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := cr.collection.FindOne(ctx, bson.M{
		"name":      name,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

func (cr *commoditiesRepository) GetByNameAndFarmerID(name string, farmerID primitive.ObjectID) (commodities.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := cr.collection.FindOne(ctx, bson.M{
		"name":      name,
		"farmerID":  farmerID,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

func (cr *commoditiesRepository) GetByQuery(query commodities.Query) ([]commodities.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"deletedAt": bson.M{"$exists": false},
	}
	if query.Name != "" {
		filter["name"] = bson.M{"$regex": query.Name, "$options": "i"}
	} else if len(query.FarmerID) != 0 {
		filter["farmerID"] = query.FarmerID
	} else if query.MinPrice != 0 && query.MaxPrice == 0 {
		filter["pricePerKg"] = bson.M{"$gte": query.MinPrice}
	} else if query.MaxPrice != 0 && query.MinPrice == 0 {
		filter["pricePerKg"] = bson.M{"$lte": query.MaxPrice}
	} else if query.MinPrice != 0 && query.MaxPrice != 0 {
		filter["pricePerKg"] = bson.M{"$gte": query.MinPrice, "$lte": query.MaxPrice}
	}

	fmt.Println(filter)

	cursor, err := cr.collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &query.Skip,
		Limit: &query.Limit,
		Sort:  bson.M{query.Sort: query.Order},
	})
	if err != nil {
		return []commodities.Domain{}, 0, err
	}

	totalData, err := cr.collection.CountDocuments(ctx, filter)
	if err != nil {
		return []commodities.Domain{}, 0, err
	}

	var result []Model
	if err = cursor.All(ctx, &result); err != nil {
		return []commodities.Domain{}, 0, err
	}

	return ToDomainArray(result), int(totalData), err
}

/*
Update
*/

func (cr *commoditiesRepository) Update(domain *commodities.Domain) (commodities.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := cr.collection.UpdateOne(ctx, bson.M{
		"_id":       domain.ID,
		"deletedAt": bson.M{"$exists": false},
	}, FromDomain(domain))
	if err != nil {
		return commodities.Domain{}, err
	}

	return *domain, err
}

/*
Delete
*/

func (cr *commoditiesRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := cr.collection.UpdateOne(ctx, bson.M{
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
