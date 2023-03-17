package transactions

import (
	"context"
	"marketplace-backend/business/transactions"
	"time"

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

/*
Update
*/

/*
Delete
*/
