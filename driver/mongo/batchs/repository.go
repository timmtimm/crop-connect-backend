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

var (
	lookupProposal = bson.M{
		"$lookup": bson.M{
			"from":         "proposals",
			"localField":   "proposalID",
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

func (br *BatchRepository) CountByProposalCode(proposalCode primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$lookup": bson.M{
				"from":         "proposals",
				"localField":   "proposalID",
				"foreignField": "_id",
				"as":           "proposal_info",
			},
		}, bson.M{
			"$project": bson.M{
				"proposal_info": bson.M{
					"$arrayElemAt": bson.A{"$proposal_info", 0},
				},
			},
		}, bson.M{
			"$match": bson.M{
				"proposal_info.code": proposalCode,
			},
		}, bson.M{
			"$count": "total",
		},
	}

	cursor, err := br.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var countResult dto.TotalDocument
	for cursor.Next(ctx) {
		err := cursor.Decode(&countResult)
		if err != nil {
			return 0, err
		}
	}

	return countResult.Total, nil
}

func (br *BatchRepository) GetByFarmerID(farmerID primitive.ObjectID) ([]batchs.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := bson.A{lookupProposal, lookupCommodity, bson.M{
		"$match": bson.M{
			"commodity_info.farmerID": farmerID,
		},
	}}
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

func (br *BatchRepository) GetByCommodityCode(commodityCode primitive.ObjectID) ([]batchs.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := bson.A{lookupProposal, lookupCommodity, bson.M{
		"$match": bson.M{
			"commodity_info.code": commodityCode,
		},
	}}
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

	if query.CommodityID != primitive.NilObjectID {
		pipeline = append(pipeline, lookupProposal, lookupCommodity, bson.M{
			"$match": bson.M{
				"proposal_info.commodityID": query.CommodityID,
			},
		})
	}

	if query.CommodityID != primitive.NilObjectID && query.FarmerID != primitive.NilObjectID {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		})
	} else if query.FarmerID != primitive.NilObjectID {
		pipeline = append(pipeline, lookupProposal, lookupCommodity, bson.M{
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

	return ToDomainArray(result), countResult.Total, nil
}

func (br *BatchRepository) CountByYear(year int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": primitive.NewDateTimeFromTime(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)),
					"$lte": primitive.NewDateTimeFromTime(time.Date(year+1, 1, 0, 0, 0, 0, 0, time.UTC)),
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

	cursor, err := br.collection.Aggregate(ctx, pipeline)
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

func (br *BatchRepository) GetForTransactionByCommodityID(commodityID primitive.ObjectID) ([]batchs.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"isAvailable": true,
			},
		}, lookupProposal, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info._id": commodityID,
			},
		},
	}

	cursor, err := br.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []batchs.Domain{}, err
	}

	var result []Model
	if err := cursor.All(ctx, &result); err != nil {
		return []batchs.Domain{}, err
	}

	return ToDomainArray(result), nil
}

func (br *BatchRepository) GetForTransactionByID(id primitive.ObjectID) (batchs.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := br.collection.FindOne(ctx, bson.M{
		"_id":         id,
		"isAvailable": true,
	}).Decode(&result)
	if err != nil {
		return batchs.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (br *BatchRepository) GetForHarvestByCommmodityID(commodityID primitive.ObjectID) ([]batchs.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$lookup": bson.M{
				"from": "harvests",
				"let":  bson.M{"batchID": "$_id"},
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{"$eq": bson.A{"$batchID", "$$batchID"}},
						}}},
				"as": "harvest_info",
			},
		}, bson.M{
			"$match": bson.M{
				"harvest_info": bson.M{"$size": 0},
			},
		}, bson.M{
			"$project": bson.M{
				"_id":                  "$_id",
				"proposalID":           "$proposalID",
				"name":                 "$name",
				"estimatedHarvestDate": "$estimatedHarvestDate",
				"status":               "$status",
				"isAvailable":          "$isAvailable",
				"createdAt":            "$createdAt",
				"harvest_info": bson.M{
					"$arrayElemAt": bson.A{"$harvest_info", 0},
				},
			},
		}, lookupProposal, lookupCommodity, bson.M{
			"$project": bson.M{
				"_id":                  "$_id",
				"proposalID":           "$proposalID",
				"name":                 "$name",
				"estimatedHarvestDate": "$estimatedHarvestDate",
				"status":               "$status",
				"isAvailable":          "$isAvailable",
				"createdAt":            "$createdAt",
				"commodity_info": bson.M{
					"$arrayElemAt": bson.A{"$commodity_info", 0},
				},
			},
		}, bson.M{
			"$match": bson.M{
				"commodity_info._id": commodityID,
			},
		},
	}

	cursor, err := br.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return []batchs.Domain{}, err
	}

	var result []Model
	if err := cursor.All(ctx, &result); err != nil {
		return []batchs.Domain{}, err
	}

	return ToDomainArray(result), nil
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
