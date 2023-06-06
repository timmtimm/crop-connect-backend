package harvests

import (
	"context"
	"crop_connect/business/harvests"
	"crop_connect/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HarvestRepository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) harvests.Repository {
	return &HarvestRepository{
		collection: db.Collection("harvests"),
	}
}

var (
	lookupBatch = bson.M{
		"$lookup": bson.M{
			"from":         "batchs",
			"localField":   "batchID",
			"foreignField": "_id",
			"as":           "batch_info",
		},
	}

	lookupProposal = bson.M{
		"$lookup": bson.M{
			"from":         "proposals",
			"localField":   "batch_info.proposalID",
			"foreignField": "_id",
			"as":           "proposal_info",
		},
	}

	lookupCommodity = bson.M{
		"$lookup": bson.M{
			"from":         "commodities",
			"localField":   "proposal_info.commodityID",
			"foreignField": "_id",
			"as":           "commodity_info",
		},
	}
)

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

func (hr *HarvestRepository) GetByID(id primitive.ObjectID) (harvests.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := hr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return harvests.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (hr *HarvestRepository) GetByBatchIDAndStatus(batchID primitive.ObjectID, status string) (harvests.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"batchID": batchID,
	}

	if status != "" {
		filter["status"] = status
	}

	var result Model
	err := hr.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return harvests.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (hr *HarvestRepository) GetByQuery(query harvests.Query) ([]harvests.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{}

	if query.Status != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"status": query.Status,
			},
		})
	}

	if query.BatchID != primitive.NilObjectID {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"batchID": query.BatchID,
			},
		})
	}

	if query.CommodityID != primitive.NilObjectID {
		pipeline = append(pipeline, lookupBatch, lookupProposal, bson.M{
			"$match": bson.M{
				"proposal_info.commodityID": query.CommodityID,
			},
		})
	}

	if query.FarmerID != primitive.NilObjectID && query.CommodityID != primitive.NilObjectID {
		pipeline = append(pipeline, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		})
	} else if query.FarmerID != primitive.NilObjectID && query.CommodityID == primitive.NilObjectID {
		pipeline = append(pipeline, lookupBatch, lookupBatch, lookupProposal, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		})
	}

	pipelineForCount := append(pipeline, bson.M{"$count": "total"})
	pipeline = append(pipeline, bson.M{
		"$skip": query.Skip,
	}, bson.M{
		"$limit": query.Limit,
	}, bson.M{
		"$sort": bson.M{query.Sort: query.Order},
	})

	cursor, err := hr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	cursorCount, err := hr.collection.Aggregate(ctx, pipelineForCount)
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

func (hr *HarvestRepository) CountByYear(year int) (float64, error) {
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
				"_id": nil,
				"total": bson.M{
					"$sum": "$totalHarvest",
				},
			},
		},
	}

	cursor, err := hr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var result []dto.TotalFloat
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	return result[0].Total, nil
}

/*
Update
*/

func (hr *HarvestRepository) Update(domain *harvests.Domain) (harvests.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := hr.collection.UpdateOne(ctx, bson.M{
		"_id": domain.ID,
	}, bson.M{
		"$set": FromDomain(domain),
	})

	if err != nil {
		return harvests.Domain{}, err
	}

	return *domain, nil
}

/*
Delete
*/
