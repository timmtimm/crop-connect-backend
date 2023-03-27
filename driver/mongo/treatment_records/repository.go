package treatment_records

import (
	"context"
	treatmentRecord "marketplace-backend/business/treatment_records"
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

/*
Update
*/

/*
Delete
*/
