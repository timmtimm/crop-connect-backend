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

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	// Update
	// Delete
}

type UseCase interface {
	// Create
	Create(domain *Domain) (int, error)
	// Read
	// Update
	// Delete
}
