package commodities

import (
	"context"
	"crop_connect/business/commodities"
	"crop_connect/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommodityRepository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) commodities.Repository {
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
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"name": bson.M{
					"$regex":   query.Name,
					"$options": "i",
				},
			},
		})
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

	if query.FarmerID != primitive.NilObjectID {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"farmerID": query.FarmerID,
			},
		})
	} else if query.Farmer != "" {
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

	if query.RegionID != primitive.NilObjectID {
		if query.Farmer == "" {
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         "users",
					"localField":   "farmerID",
					"foreignField": "_id",
					"as":           "farmer_info",
				},
			})
		}

		match := bson.M{
			"$match": bson.M{
				"farmer_info.regionID": query.RegionID,
			},
		}

		pipeline = append(pipeline, match)
	} else if query.Province != "" || query.Regency != "" || query.District != "" {
		if query.Farmer == "" {
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         "users",
					"localField":   "farmerID",
					"foreignField": "_id",
					"as":           "farmer_info",
				},
			})
		}

		lookup := bson.M{
			"$lookup": bson.M{
				"from":         "regions",
				"localField":   "farmer_info.regionID",
				"foreignField": "_id",
				"as":           "region_info",
			},
		}

		pipeline = append(pipeline, lookup)

		if query.Province != "" {
			match := bson.M{
				"$match": bson.M{
					"region_info.province": query.Province,
				},
			}

			pipeline = append(pipeline, match)
		}

		if query.Regency != "" {
			match := bson.M{
				"$match": bson.M{
					"region_info.regency": query.Regency,
				},
			}

			pipeline = append(pipeline, match)
		}

		if query.District != "" {
			match := bson.M{
				"$match": bson.M{
					"region_info.district": query.District,
				},
			}

			pipeline = append(pipeline, match)
		}
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

	pipelineForCount := make([]interface{}, len(pipeline))
	copy(pipelineForCount, pipeline)
	pipelineForCount = append(pipelineForCount, bson.M{
		"$count": "total",
	})

	// Convert to JSON
	// jsonData, err := json.Marshal(pipelineForCount)
	// if err != nil {
	// 	panic(err)
	// }

	// // Print JSON
	// fmt.Println("SEBELUM APPEND")
	// fmt.Println(string(jsonData))

	pipeline = append(pipeline, paginationSkip, paginationLimit, paginationSort)

	cursor, err := cr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	// Convert to JSON
	// jsonData, _ = json.Marshal(pipelineForCount)

	// Print JSON
	// fmt.Println("SETELAH APPEND")
	// fmt.Println(string(jsonData))

	cursorCount, err := cr.collection.Aggregate(ctx, pipelineForCount)
	if err != nil {
		return nil, 0, err
	}

	var result []Model
	var countResult dto.TotalDocument

	if err := cursor.All(ctx, &result); err != nil {
		return nil, 0, err
	}

	for cursorCount.Next(ctx) {
		err := cursorCount.Decode(&countResult)
		if err != nil {
			return nil, 0, err
		}
	}

	return ToDomainArray(result), countResult.Total, nil
}

func (cr *CommodityRepository) CountTotalCommodity(year int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": primitive.NewDateTimeFromTime(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)),
					"$lte": primitive.NewDateTimeFromTime(time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				"deletedAt": bson.M{"$exists": false},
			},
		}, bson.M{
			"$group": bson.M{
				"_id": year,
				"total": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	cursor, err := cr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var result dto.TotalDocument
	for cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return 0, err
		}
	}

	return result.Total, nil
}

func (cr *CommodityRepository) CountTotalCommodityByFarmer(farmerID primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	count, err := cr.collection.CountDocuments(ctx, bson.M{
		"farmerID": farmerID,
	})

	return int(count), err
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
