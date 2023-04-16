package harvests

import (
	"context"
	"marketplace-backend/business/harvests"
	"marketplace-backend/dto"
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

func (hr *HarvestRepository) GetByBatchID(batchID primitive.ObjectID) (harvests.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := hr.collection.FindOne(ctx, bson.M{
		"batchID": batchID,
	}).Decode(&result)
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

	lookupBatch := bson.M{
		"$lookup": bson.M{
			"from":         "batchs",
			"localField":   "batchID",
			"foreignField": "_id",
			"as":           "batch_info",
		},
	}

	lookupTransaction := bson.M{
		"$lookup": bson.M{
			"from":         "transactions",
			"localField":   "batch_info.transactionID",
			"foreignField": "_id",
			"as":           "transaction_info",
		},
	}

	lookupProposal := bson.M{
		"$lookup": bson.M{
			"from":         "proposals",
			"localField":   "transaction_info.proposalID",
			"foreignField": "_id",
			"as":           "proposal_info",
		},
	}

	lookupCommodity := bson.M{
		"$lookup": bson.M{
			"from":         "commodities",
			"localField":   "proposal_info.commodityID",
			"foreignField": "_id",
			"as":           "commodity_info",
		},
	}

	if query.Batch != "" {
		pipeline = append(pipeline, lookupBatch, bson.M{
			"$match": bson.M{
				"batch_info.name": query.Batch,
			},
		})
	}

	if query.Commodity != "" && query.Batch != "" {
		pipeline = append(pipeline, lookupTransaction, lookupProposal, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info.name": query.Commodity,
			},
		})
	} else if query.Commodity != "" && query.Batch == "" {
		pipeline = append(pipeline, lookupBatch, lookupTransaction, lookupProposal, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info.name": query.Commodity,
			},
		})
	}

	if query.FarmerID != primitive.NilObjectID && query.Commodity != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		})
	} else if query.FarmerID != primitive.NilObjectID && query.Commodity == "" {
		pipeline = append(pipeline, lookupBatch, lookupTransaction, lookupProposal, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		})
	}

	pipelineForCount := append(pipeline, bson.M{"$count": "totalDocument"})
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

	return ToDomainArray(result), countResult.TotalDocument, nil
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
