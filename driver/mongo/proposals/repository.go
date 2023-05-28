package proposals

import (
	"context"
	"crop_connect/business/proposals"
	"crop_connect/constant"
	"crop_connect/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProposalRepository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) proposals.Repository {
	return &ProposalRepository{
		collection: db.Collection("proposals"),
	}
}

var (
	lookupCommodity = bson.M{
		"$lookup": bson.M{
			"from":         "commodities",
			"localField":   "commodityID",
			"foreignField": "_id",
			"as":           "commodity_info",
		},
	}
)

/*
Create
*/

func (pr *ProposalRepository) Create(domain *proposals.Domain) (proposals.Domain, error) {
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

func (pr *ProposalRepository) GetByID(id primitive.ObjectID) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := pr.collection.FindOne(ctx, bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

func (pr *ProposalRepository) GetByIDWithoutDeleted(id primitive.ObjectID) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := pr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (pr *ProposalRepository) GetByCommodityID(commodityID primitive.ObjectID) ([]proposals.Domain, error) {
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

func (pr *ProposalRepository) GetByCommodityIDAndAvailability(commodityID primitive.ObjectID, status string) ([]proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := pr.collection.Find(ctx, bson.M{
		"commodityID": commodityID,
		"status":      status,
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

func (pr *ProposalRepository) GetByCommodityIDAndName(commodityID primitive.ObjectID, name string) (proposals.Domain, error) {
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

func (pr *ProposalRepository) GetByIDAccepted(id primitive.ObjectID) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := pr.collection.FindOne(ctx, bson.M{
		"_id":       id,
		"status":    constant.ProposalStatusApproved,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

func (pr *ProposalRepository) StatisticByYear(year int) ([]dto.StatisticByYear, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"status": constant.ProposalStatusApproved,
				"createdAt": bson.M{
					"$gte": time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
					"$lte": time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				"deletedAt": bson.M{"$exists": false},
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

	var result []dto.StatisticByYear
	cursor, err := pr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []dto.StatisticByYear{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []dto.StatisticByYear{}, err
	}

	return result, err
}

func (pr *ProposalRepository) CountTotalProposalByFarmer(farmerID primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"deletedAt": bson.M{"$exists": false},
			},
		}, lookupCommodity, bson.M{
			"$project": bson.M{
				"commodity_info": bson.M{
					"$arrayElemAt": bson.A{"$commodity_info", 0},
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

	var result dto.TotalDocument
	cursor, err := pr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return 0, err
		}
	}

	return result.Total, nil
}

func (pr *ProposalRepository) GetByQuery(query proposals.Query) ([]proposals.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"deletedAt": bson.M{"$exists": false},
			},
		},
	}

	if query.CommodityID != primitive.NilObjectID {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"commodityID": query.CommodityID,
			},
		})
	}

	if query.FarmerID != primitive.NilObjectID {
		pipeline = append(pipeline, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		})
	}

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

	if query.Status != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"status": query.Status,
			},
		})
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

	pipeline = append(pipeline, paginationSkip, paginationLimit, paginationSort)

	cursor, err := pr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []proposals.Domain{}, 0, err
	}

	cursorCount, err := pr.collection.Aggregate(ctx, pipelineForCount)
	if err != nil {
		return []proposals.Domain{}, 0, err
	}

	var result []Model
	var countResult dto.TotalDocument

	if err := cursor.All(ctx, &result); err != nil {
		return []proposals.Domain{}, 0, err
	}

	for cursorCount.Next(ctx) {
		err := cursorCount.Decode(&countResult)
		if err != nil {
			return []proposals.Domain{}, 0, err
		}
	}

	return ToDomainArray(result), countResult.Total, nil
}

/*
Update
*/

func (pr *ProposalRepository) Update(domain *proposals.Domain) (proposals.Domain, error) {
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

func (pr *ProposalRepository) UnsetRejectReason(id primitive.ObjectID) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := pr.collection.UpdateOne(ctx, bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}, bson.M{
		"$unset": bson.M{"rejectReason": ""},
	})
	if err != nil {
		return proposals.Domain{}, err
	}

	updatedProposal, err := pr.GetByID(id)
	if err != nil {
		return proposals.Domain{}, err
	}

	return updatedProposal, nil
}

/*
Delete
*/

func (pr *ProposalRepository) Delete(id primitive.ObjectID) error {
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

func (pr *ProposalRepository) DeleteByCommodityID(commodityID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := pr.collection.UpdateMany(ctx, bson.M{
		"commodityID": commodityID,
		"deletedAt":   bson.M{"$exists": false},
	}, bson.M{
		"$set": bson.M{
			"deletedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	})

	return err
}
