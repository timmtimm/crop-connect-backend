package transactions

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	ID         primitive.ObjectID
	BuyerID    primitive.ObjectID
	ProposalID primitive.ObjectID
	Address    string
	Status     string
	TotalPrice float64
	CreatedAt  primitive.DateTime
	UpdatedAt  primitive.DateTime
}

type Query struct {
	Skip      int64
	Limit     int64
	Sort      string
	Order     int
	Commodity string
	FarmerID  primitive.ObjectID
	BuyerID   primitive.ObjectID
	Status    string
	StartDate primitive.DateTime
	EndDate   primitive.DateTime
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByBuyerIDProposalIDAndStatus(buyerID primitive.ObjectID, proposalID primitive.ObjectID, status string) (Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	// Update
	// Delete
}

type UseCase interface {
	// Create
	Create(domain *Domain) (int, error)
	// Read
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	// Update
	// Delete
}
