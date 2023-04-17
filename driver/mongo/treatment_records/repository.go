package treatment_records

import (
	"context"
	treatmentRecord "crop_connect/business/treatment_records"
	"crop_connect/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TreatmentRecordRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) treatmentRecord.Repository {
	return &TreatmentRecordRepository{
		collection: db.Collection("treatmentRecords"),
	}
}

/*
Create
*/

func (trr *TreatmentRecordRepository) Create(domain *treatmentRecord.Domain) (treatmentRecord.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := trr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return treatmentRecord.Domain{}, err
	}

	return *domain, err
}

/*
Read
*/

func (trr *TreatmentRecordRepository) GetNewestByBatchID(batchID primitive.ObjectID) (treatmentRecord.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := trr.collection.FindOne(ctx, bson.M{
		"batchID": batchID,
	}, &options.FindOneOptions{
		Sort: bson.M{
			"createdAt": -1,
		},
	}).Decode(&result)
	if err != nil {
		return treatmentRecord.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (trr *TreatmentRecordRepository) CountByBatchID(batchID primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	count, err := trr.collection.CountDocuments(ctx, bson.M{
		"batchID": batchID,
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (trr *TreatmentRecordRepository) GetByID(id primitive.ObjectID) (treatmentRecord.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := trr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (trr *TreatmentRecordRepository) GetByBatchID(batchID primitive.ObjectID) ([]treatmentRecord.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := trr.collection.Find(ctx, bson.M{
		"batchID": batchID,
	})
	if err != nil {
		return []treatmentRecord.Domain{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []treatmentRecord.Domain{}, err
	}

	return ToDomainArray(result), nil
}

func (trr *TreatmentRecordRepository) GetByQuery(query treatmentRecord.Query) ([]treatmentRecord.Domain, int, error) {
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

	if query.Number != 0 {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"number": query.Number,
			},
		})
	}

	if query.Batch != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": query.Batch,
				},
			},
		})
	}

	if query.Commodity != "" {
		lookup1 := bson.M{
			"$lookup": bson.M{
				"from":         "batchs",
				"localField":   "batchID",
				"foreignField": "_id",
				"as":           "batch_info",
			},
		}

		lookup2 := bson.M{
			"$lookup": bson.M{
				"from":         "transactions",
				"localField":   "batch_info.transactionID",
				"foreignField": "_id",
				"as":           "transaction_info",
			},
		}

		lookup3 := bson.M{
			"$lookup": bson.M{
				"from":         "proposals",
				"localField":   "transaction_info.proposalID",
				"foreignField": "_id",
				"as":           "proposal_info",
			},
		}

		lookup4 := bson.M{
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

		pipeline = append(pipeline, lookup1, lookup2, lookup3, lookup4, match)
	}

	if query.FarmerID != primitive.NilObjectID && query.Commodity == "" {
		lookup1 := bson.M{
			"$lookup": bson.M{
				"from":         "batchs",
				"localField":   "batchID",
				"foreignField": "_id",
				"as":           "batch_info",
			},
		}

		lookup2 := bson.M{
			"$lookup": bson.M{
				"from":         "transactions",
				"localField":   "batch_info.transactionID",
				"foreignField": "_id",
				"as":           "transaction_info",
			},
		}

		lookup3 := bson.M{
			"$lookup": bson.M{
				"from":         "proposals",
				"localField":   "transaction_info.proposalID",
				"foreignField": "_id",
				"as":           "proposal_info",
			},
		}

		lookup4 := bson.M{
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

		pipeline = append(pipeline, lookup1, lookup2, lookup3, lookup4, match)
	} else if query.FarmerID != primitive.NilObjectID && query.Commodity != "" {
		match := bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		}

		pipeline = append(pipeline, match)
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

	cursor, err := trr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	cursorCount, err := trr.collection.Aggregate(ctx, pipelineForCount)
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

func (trr *TreatmentRecordRepository) Update(domain *treatmentRecord.Domain) (treatmentRecord.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := trr.collection.UpdateOne(ctx, bson.M{
		"_id": domain.ID,
	}, bson.M{
		"$set": FromDomain(domain),
	})

	if err != nil {
		return treatmentRecord.Domain{}, err
	}

	return *domain, nil
}

/*
Delete
*/
