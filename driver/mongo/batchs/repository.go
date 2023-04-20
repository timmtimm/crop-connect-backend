package batchs

import (
	"context"
	"crop_connect/business/batchs"
	"crop_connect/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BatchRepository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) batchs.Repository {
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

func (br *BatchRepository) GetByID(id primitive.ObjectID) (batchs.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := br.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)

	return result.ToDomain(), err
}

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

func (br *BatchRepository) GetByFarmerID(farmerID primitive.ObjectID) ([]batchs.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	lookupTransaction := bson.M{
		"$lookup": bson.M{
			"from":         "transactions",
			"localField":   "transactionID",
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

	match := bson.M{
		"$match": bson.M{
			"commodity_info.farmerID": farmerID,
		},
	}

	pipeline := bson.A{lookupTransaction, lookupProposal, lookupCommodity, match}
	cursor, err := br.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []Model
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return ToDomainArray(result), nil
}

func (br *BatchRepository) GetByCommodityID(commodityID primitive.ObjectID) ([]batchs.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	lookupTransaction := bson.M{
		"$lookup": bson.M{
			"from":         "transactions",
			"localField":   "transactionID",
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

	match := bson.M{
		"$match": bson.M{
			"proposal_info.commodityID": commodityID,
		},
	}

	pipeline := bson.A{lookupTransaction, lookupProposal, match}
	cursor, err := br.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []Model
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return ToDomainArray(result), nil
}

func (br *BatchRepository) GetByQuery(query batchs.Query) ([]batchs.Domain, int, error) {
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

	lookupTransaction := bson.M{
		"$lookup": bson.M{
			"from":         "transactions",
			"localField":   "transactionID",
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

	if query.Commodity != "" {
		pipeline = append(pipeline, lookupTransaction, lookupProposal, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info.name": bson.M{
					"$regex":   query.Commodity,
					"$options": "i",
				},
			},
		})
	}

	if query.Commodity != "" && query.FarmerID != primitive.NilObjectID {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		})
	} else if query.FarmerID != primitive.NilObjectID {
		pipeline = append(pipeline, lookupTransaction, lookupProposal, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		})
	}

	cursor, err := br.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	cursorCount, err := br.collection.Aggregate(ctx, append(pipeline, bson.M{"$count": "totalDocument"}))
	if err != nil {
		return nil, 0, err
	}

	countResult := dto.TotalDocument{}
	for cursorCount.Next(ctx) {
		err := cursorCount.Decode(&countResult)
		if err != nil {
			return nil, 0, err
		}
	}

	var result []Model
	if err := cursor.All(ctx, &result); err != nil {
		return nil, 0, err
	}

	return ToDomainArray(result), countResult.TotalDocument, nil
}

/*
Update
*/

func (br *BatchRepository) Update(domain *batchs.Domain) (batchs.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := br.collection.UpdateOne(ctx, bson.M{
		"_id": domain.ID,
	}, bson.M{
		"$set": FromDomain(domain),
	})

	if err != nil {
		return batchs.Domain{}, err
	}

	return *domain, nil
}

/*
Delete
*/
