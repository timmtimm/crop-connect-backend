package commodities

import (
	"context"
	"marketplace-backend/business/commodities"
	"marketplace-backend/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommodityRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) commodities.Repository {
	return &CommodityRepository{
		collection: db.Collection("commodities"),
	}
}

/*
Create
*/

func (cr *CommodityRepository) Create(domain *commodities.Domain) (commodities.Domain, error) {
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

func (cr *CommodityRepository) GetByID(id primitive.ObjectID) (commodities.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := cr.collection.FindOne(ctx, bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

func (cr *CommodityRepository) GetByIDWithoutDeleted(id primitive.ObjectID) (commodities.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := cr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (cr *CommodityRepository) GetByIDAndFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (commodities.Domain, error) {
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

func (cr *CommodityRepository) GetByName(name string) (commodities.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := cr.collection.FindOne(ctx, bson.M{
		"name":      name,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

func (cr *CommodityRepository) GetByNameAndFarmerID(name string, farmerID primitive.ObjectID) (commodities.Domain, error) {
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

func (cr *CommodityRepository) GetByFarmerID(farmerID primitive.ObjectID) ([]commodities.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := cr.collection.Find(ctx, bson.M{
		"farmerID":  farmerID,
		"deletedAt": bson.M{"$exists": false},
	})
	if err != nil {
		return []commodities.Domain{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []commodities.Domain{}, err
	}

	return ToDomainArray(result), err
}

func (cr *CommodityRepository) GetByQuery(query commodities.Query) ([]commodities.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{}
	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"deletedAt": bson.M{"$exists": false},
		},
	})

	if query.Name != "" {
		filterName := bson.M{"$regex": query.Name, "$options": "i"}
		pipeline = append(pipeline, filterName)
	}

	if query.FarmerID != primitive.NilObjectID {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"farmerID": query.FarmerID,
			},
		})
	}

	if query.MinPrice != 0 {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"pricePerKg": bson.M{
					"$gte": query.MinPrice,
				},
			},
		})
	}

	if query.MaxPrice != 0 {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"pricePerKg": bson.M{
					"$lte": query.MaxPrice,
				},
			},
		})
	}

	if query.Farmer != "" {
		lookup := bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "farmerID",
				"foreignField": "_id",
				"as":           "farmer_info",
			},
		}

		match := bson.M{
			"$match": bson.M{
				"farmer_info.name": bson.M{
					"$regex":   query.Farmer,
					"$options": "i",
				},
			},
		}

		pipeline = append(pipeline, lookup, match)
	}

	paginationSkip := bson.M{
		"$skip": query.Skip,
	}

	paginationLimit := bson.M{
		"$limit": query.Limit,
	}

	paginationSort := bson.M{
		"$sort": bson.M{query.Sort: query.Order},
	}

	pipelineForCount := append(pipeline, bson.M{"$count": "totalDocument"})
	pipeline = append(pipeline, paginationSkip, paginationLimit, paginationSort)

	cursor, err := cr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	cursorCount, err := cr.collection.Aggregate(ctx, pipelineForCount)
	if err != nil {
		return nil, 0, err
	}

	var result []Model
	countResult := dto.TotalDocument{}

	if err := cursor.All(ctx, &result); err != nil {
		return nil, 0, err
	}

	for cursorCount.Next(ctx) {
		err := cursorCount.Decode(&countResult)
		if err != nil {
			return nil, 0, err
		}

	}

	return ToDomainArray(result), countResult.TotalDocument, nil
}

/*
Update
*/

func (cr *CommodityRepository) Update(domain *commodities.Domain) (commodities.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := cr.collection.UpdateOne(ctx, bson.M{
		"_id":       domain.ID,
		"deletedAt": bson.M{"$exists": false},
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return commodities.Domain{}, err
	}

	return *domain, err
}

/*
Delete
*/

func (cr *CommodityRepository) Delete(id primitive.ObjectID) error {
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
