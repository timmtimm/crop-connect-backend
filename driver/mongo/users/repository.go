package users

import (
	"context"
	"crop_connect/business/users"
	"crop_connect/dto"
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

	pipeline := []interface{}{}

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

	if query.Email != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"email": bson.M{
					"$regex":   query.Email,
					"$options": "i",
				},
			},
		})
	}

	if query.PhoneNumber != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"phoneNumber": bson.M{
					"$regex": query.PhoneNumber,
				},
			},
		})
	}

	if query.Role != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"role": query.Role,
			},
		})
	}

	if query.RegionID != primitive.NilObjectID {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"regionID": query.RegionID,
			},
		})
	} else if query.Province != "" || query.Regency != "" || query.District != "" {
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from":         "regions",
				"localField":   "regionID",
				"foreignField": "_id",
				"as":           "region_info",
			},
		})

		if query.Province != "" {
			pipeline = append(pipeline, bson.M{
				"$match": bson.M{
					"region_info.province": query.Province,
				},
			})
		}

		if query.Regency != "" {
			pipeline = append(pipeline, bson.M{
				"$match": bson.M{
					"region_info.regency": query.Regency,
				},
			})
		}

		if query.District != "" {
			pipeline = append(pipeline, bson.M{
				"$match": bson.M{
					"region_info.district": query.District,
				},
			})
		}
	}

	pipelineForCount := append(pipeline, bson.M{"$count": "total"})
	pipeline = append(pipeline, bson.M{
		"$skip": query.Skip,
	}, bson.M{
		"$limit": query.Limit,
	}, bson.M{
		"$sort": bson.M{query.Sort: query.Order},
	})

	cursor, err := ur.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	cursorCount, err := ur.collection.Aggregate(ctx, pipelineForCount)
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

	return ToDomainArray(result), countResult.Total, nil
}

func (ur *UserRepository) GetFarmerByID(id primitive.ObjectID) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ur.collection.FindOne(ctx, bson.M{
		"_id":  id,
		"role": "farmer",
	}).Decode(&result)

	return result.ToDomain(), err
}

func (ur *UserRepository) StatisticNewUserByYear(year int) ([]dto.StatisticByYear, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
					"$lte": time.Date(year+1, 0, 0, 0, 0, 0, 0, time.UTC),
				},
			},
		}, bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"$month": "$createdAt",
				},
				"total": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	cursor, err := ur.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []dto.StatisticByYear{}, err
	}

	var result []dto.StatisticByYear
	if err := cursor.All(ctx, &result); err != nil {
		return []dto.StatisticByYear{}, err
	}

	return result, nil
}

func (ur *UserRepository) CountTotalValidatorByYear(year int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"role": "validator",
				"createdAt": bson.M{
					"$gte": time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
					"$lte": time.Date(year+1, 0, 0, 0, 0, 0, 0, time.UTC),
				},
			},
		}, bson.M{
			"$group": bson.M{
				"_id": nil,
				"total": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	cursor, err := ur.collection.Aggregate(ctx, pipeline)
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
