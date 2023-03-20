package transactions

import (
	"context"
	"marketplace-backend/business/transactions"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type transactionRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) transactions.Repository {
	return &transactionRepository{
		collection: db.Collection("transactions"),
	}
}

/*
Create
*/

func (tr *transactionRepository) Create(domain *transactions.Domain) (transactions.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return transactions.Domain{}, err
	}

	return *domain, err
}

/*
Read
*/

func (tr *transactionRepository) GetByBuyerIDProposalIDAndStatus(buyerID primitive.ObjectID, proposalID primitive.ObjectID, status string) (transactions.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := tr.collection.FindOne(ctx, bson.M{
		"buyerID":    buyerID,
		"proposalID": proposalID,
		"status":     status,
		"deletedAt":  bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

func (tr *transactionRepository) GetByQuery(query transactions.Query) ([]transactions.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{}
	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"deletedAt": bson.M{"$exists": false},
		},
	})

	if query.BuyerID != primitive.NilObjectID {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"buyerID": query.BuyerID,
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

	if query.StartDate != 0 {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": query.StartDate,
				},
			},
		})
	}

	if query.EndDate != 0 {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$lte": query.EndDate,
				},
			},
		})
	}

	if query.Commodity != "" {
		lookup1 := bson.M{
			"$lookup": bson.M{
				"from":         "proposals",
				"localField":   "proposalID",
				"foreignField": "_id",
				"as":           "proposal_info",
			},
		}

		lookup2 := bson.M{
			"$lookup": bson.M{
				"from":         "commodities",
				"localField":   "proposal_info.commodityID",
				"foreignField": "_id",
				"as":           "commodity_info",
			},
		}

		match := bson.M{
			"$match": bson.M{
				"commodity_info.name": bson.M{
					"$regex":   query.Commodity,
					"$options": "i",
				},
			},
		}

		pipeline = append(pipeline, lookup1, lookup2, match)
	}

	if query.FarmerID != primitive.NilObjectID {
		lookup1 := bson.M{
			"$lookup": bson.M{
				"from":         "proposals",
				"localField":   "proposalID",
				"foreignField": "_id",
				"as":           "proposal_info",
			},
		}

		lookup2 := bson.M{
			"$lookup": bson.M{
				"from":         "commodities",
				"localField":   "proposal_info.commodityID",
				"foreignField": "_id",
				"as":           "commodity_info",
			},
		}

		match := bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		}

		pipeline = append(pipeline, lookup1, lookup2, match)
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

	cursor, err := tr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	cursorCount, err := tr.collection.Aggregate(ctx, pipelineForCount)
	if err != nil {
		return nil, 0, err
	}

	var result []Model
	countResult := TotalDocument{}

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

/*
Delete
*/
