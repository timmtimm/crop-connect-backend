package batchs

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	ID                   primitive.ObjectID
	ProposalID           primitive.ObjectID
	Name                 string
	EstimatedHarvestDate primitive.DateTime
	Status               string
	CancelReason         string
	IsAvailable          bool
	CreatedAt            primitive.DateTime
	UpdatedAt            primitive.DateTime
}

type Query struct {
	Skip        int64
	Limit       int64
	Sort        string
	Order       int
	FarmerID    primitive.ObjectID
	CommodityID primitive.ObjectID
	Name        string
	Status      string
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	CountByProposalCode(proposalCode primitive.ObjectID) (int, error)
	GetByFarmerID(farmerID primitive.ObjectID) ([]Domain, error)
	GetByCommodityCode(commodityCode primitive.ObjectID) ([]Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	CountByYear(year int) (int, error)
	GetForTransactionByCommodityID(commodityID primitive.ObjectID) ([]Domain, error)
	GetForTransactionByCommodityCode(commodityCode primitive.ObjectID) ([]Domain, error)
	GetForTransactionByID(id primitive.ObjectID) (Domain, error)
	GetForHarvestByCommmodityID(commodityID primitive.ObjectID) ([]Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
}

type UseCase interface {
	// Create
	Create(proposalID primitive.ObjectID, farmerID primitive.ObjectID) (int, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, int, error)
	GetByCommodityID(commodityID primitive.ObjectID) ([]Domain, int, error)
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	CountByYear(year int) (int, int, error)
	GetForTransactionByCommodityID(commodityID primitive.ObjectID) ([]Domain, int, error)
	GetForTransactionByID(id primitive.ObjectID) (Domain, int, error)
	GetForHarvestByCommmodityIDAndFarmerID(commodityID primitive.ObjectID, farmerID primitive.ObjectID) ([]Domain, int, error)
	// Update
	// Cancel(domain *Domain, farmerID primitive.ObjectID) (int, error)
	// Delete
}
