package transactions

import (
	"context"
	"crop_connect/business/transactions"
	"crop_connect/constant"
	"crop_connect/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) transactions.Repository {
	return &TransactionRepository{
		collection: db.Collection("transactions"),
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

func (tr *TransactionRepository) Create(domain *transactions.Domain) (transactions.Domain, error) {
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

func (tr *TransactionRepository) GetByID(id primitive.ObjectID) (transactions.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := tr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (tr *TransactionRepository) GetByBuyerIDProposalIDAndStatus(buyerID primitive.ObjectID, proposalID primitive.ObjectID, status string) (transactions.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := tr.collection.FindOne(ctx, bson.M{
		"buyerID":    buyerID,
		"proposalID": proposalID,
		"status":     status,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (tr *TransactionRepository) GetByQuery(query transactions.Query) ([]transactions.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{}

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
		pipeline = append(pipeline, lookupProposal, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info.name": bson.M{
					"$regex":   query.Commodity,
					"$options": "i",
				},
			},
		})
	}

	if query.FarmerID != primitive.NilObjectID && query.Commodity == "" {
		pipeline = append(pipeline, lookupProposal, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		})
	} else if query.FarmerID != primitive.NilObjectID && query.Commodity != "" {
		pipeline = append(pipeline, bson.M{
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

	cursor, err := tr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	cursorCount, err := tr.collection.Aggregate(ctx, pipelineForCount)
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

func (tr *TransactionRepository) GetByIDAndBuyerID(id primitive.ObjectID, buyerID primitive.ObjectID) (transactions.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := tr.collection.FindOne(ctx, bson.M{
		"_id":     id,
		"buyerID": buyerID,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (tr *TransactionRepository) StatisticByYear(farmerID primitive.ObjectID, year int) ([]transactions.Statistic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var results []transactions.Statistic
	var pipeline []interface{}

	for month := 1; month <= 12; month++ {
		if time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Unix() > time.Now().Unix() {
			results = append(results, transactions.Statistic{
				Month:            int(time.Month(month)),
				TotalAccepted:    0,
				TotalTransaction: 0,
				TotalIncome:      0,
				TotalWeight:      0,
				TotalUniqueBuyer: 0,
			})

			continue
		}

		if month < 12 {
			pipeline = []interface{}{
				bson.M{
					"$match": bson.M{
						"createdAt": bson.M{
							"$gte": primitive.NewDateTimeFromTime(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)),
							"$lte": primitive.NewDateTimeFromTime(time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
			}
		} else {
			pipeline = []interface{}{
				bson.M{
					"$match": bson.M{
						"createdAt": bson.M{
							"$gte": primitive.NewDateTimeFromTime(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)),
							"$lte": primitive.NewDateTimeFromTime(time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)),
						},
					},
				},
			}
		}

		if farmerID != primitive.NilObjectID {
			pipeline = append(pipeline, lookupProposal, lookupCommodity, bson.M{
				"$match": bson.M{
					"commodity_info.farmerID": farmerID,
				},
			})
		}

		pipeline = append(pipeline,
			bson.M{
				"$lookup": bson.M{
					"from":         "batchs",
					"localField":   "_id",
					"foreignField": "transactionID",
					"as":           "batch_info",
				},
			}, bson.M{
				"$lookup": bson.M{
					"from":         "harvests",
					"localField":   "batch_info._id",
					"foreignField": "batchID",
					"as":           "harvest_info",
				},
			}, bson.M{
				"$project": bson.M{
					"buyerID":    "$buyerID",
					"status":     "$status",
					"totalPrice": "$totalPrice",
					"harvest_info": bson.M{
						"$arrayElemAt": bson.A{"$harvest_info", 0},
					},
				},
			}, bson.M{
				"$group": bson.M{
					"_id": "transaction_statistic",
					"totalAccepted": bson.M{
						"$sum": bson.M{
							"$cond": bson.A{
								bson.M{
									"$eq": bson.A{
										"$status", constant.TransactionStatusAccepted,
									},
								}, 1, 0,
							},
						},
					},
					"totalTransaction": bson.M{
						"$sum": 1,
					},
					"totalIncome": bson.M{
						"$sum": "$totalPrice",
					},
					"totalWeight": bson.M{
						"$sum": bson.M{
							"$cond": bson.A{
								bson.M{"$eq": bson.A{"$harvest_info.status", constant.HarvestStatusApproved}},
								"$harvest_info.totalHarvest", 0},
						},
					},
					"uniqueBuyer": bson.M{
						"$addToSet": "$buyerID",
					},
				},
			}, bson.M{
				"$project": bson.M{
					"totalAccepted":    "$totalAccepted",
					"totalTransaction": "$totalTransaction",
					"totalIncome":      "$totalIncome",
					"totalWeight":      "$totalWeight",
					"totalUniqueBuyer": bson.M{
						"$size": "$uniqueBuyer",
					},
				},
			},
		)

		cursor, err := tr.collection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		var result StatisticModel
		for cursor.Next(ctx) {
			err := cursor.Decode(&result)
			if err != nil {
				return nil, err
			}
		}
		result.Month = month

		results = append(results, result.ToStatistic())
	}

	return results, nil
}

func (tr *TransactionRepository) StatisticTopProvince(year int, limit int) ([]transactions.TotalTransactionByProvince, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": primitive.NewDateTimeFromTime(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)),
					"$lte": primitive.NewDateTimeFromTime(time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		}, bson.M{
			"$lookup": bson.M{
				"from":         "regions",
				"localField":   "regionID",
				"foreignField": "_id",
				"as":           "region_info",
			},
		}, bson.M{
			"$project": bson.M{
				"status": "$status",
				"region_info": bson.M{
					"$arrayElemAt": bson.A{"$region_info", 0},
				},
			},
		}, bson.M{
			"$group": bson.M{
				"_id": "$region_info.province",
				"totalAccepted": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$status", constant.TransactionStatusAccepted}},
							1, 0},
					},
				},
				"totalTransaction": bson.M{
					"$sum": 1,
				},
			},
		}, bson.M{
			"$sort": bson.M{
				"total": -1,
			},
		}, bson.M{
			"$limit": limit,
		},
	}

	cursor, err := tr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []TotalTransactionByProvince
	if err := cursor.All(ctx, &results); err != nil {
		return []transactions.TotalTransactionByProvince{}, err
	}

	return ToTotalTransactionByProvinceArray(results), nil
}

func (tr *TransactionRepository) StatisticTopCommodity(farmerID primitive.ObjectID, year int, limit int) ([]transactions.ModelStatisticTopCommodity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": primitive.NewDateTimeFromTime(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)),
					"$lte": primitive.NewDateTimeFromTime(time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		}, lookupProposal, lookupCommodity,
	}

	if farmerID != primitive.NilObjectID {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": farmerID,
			},
		})
	}

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": "$commodity_info._id",
			"commodity_info": bson.M{
				"$arrayElemAt": bson.A{"$commodity_info", 0},
			},
		},
	}, bson.M{
		"$group": bson.M{
			"_id": "$commodity_info._id",
			"total": bson.M{
				"$sum": 1,
			},
		},
	}, bson.M{
		"$sort": bson.M{
			"total": -1,
		},
	}, bson.M{
		"$limit": limit,
	})

	cursor, err := tr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []transactions.ModelStatisticTopCommodity
	if err := cursor.All(ctx, &results); err != nil {
		return []transactions.ModelStatisticTopCommodity{}, err
	}

	return results, nil
}

/*
Update
*/

func (tr *TransactionRepository) Update(domain *transactions.Domain) (transactions.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.UpdateOne(ctx, bson.M{
		"_id": domain.ID,
	}, bson.M{
		"$set": FromDomain(domain),
	})

	if err != nil {
		return transactions.Domain{}, err
	}

	return *domain, nil
}

func (tr *TransactionRepository) RejectPendingByProposalID(proposalID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.UpdateMany(ctx, bson.M{
		"proposalID": proposalID,
		"status":     constant.TransactionStatusPending,
	}, bson.M{
		"$set": bson.M{
			"status":    constant.TransactionStatusRejected,
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	})

	return err
}

/*
Delete
*/
