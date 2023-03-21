package batchs

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	ID                   primitive.ObjectID
	TransactionID        primitive.ObjectID
	Name                 string
	EstimatedHarvestDate primitive.DateTime
	Status               string
	CancelReason         string
	CreatedAt            primitive.DateTime
	UpdatedAt            primitive.DateTime
}

type Query struct {
	Skip      int64
	Limit     int64
	Sort      string
	Order     int
	Commodity string
	Status    string
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	CountByProposalName(proposalName string) (int, error)
	// Update
	// Delete
}

type UseCase interface {
	// Create
	Create(transactionID primitive.ObjectID) (int, error)
	// Read
	// Update
	// Delete
}
