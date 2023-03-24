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
	FarmerID  primitive.ObjectID
	Commodity string
	Name      string
	Status    string
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	CountByProposalName(proposalName string) (int, error)
	GetByFarmerID(farmerID primitive.ObjectID) ([]Domain, error)
	GetByCommodityID(commodityID primitive.ObjectID) ([]Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
}

type UseCase interface {
	// Create
	Create(transactionID primitive.ObjectID) (int, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, int, error)
	GetByCommodityID(commodityID primitive.ObjectID) ([]Domain, int, error)
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	// Update
	Cancel(domain *Domain, farmerID primitive.ObjectID) (int, error)
	// Delete
}
