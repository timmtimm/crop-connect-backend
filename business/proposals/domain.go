package proposals

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	ID                    primitive.ObjectID
	AccepterID            primitive.ObjectID
	CommodityID           primitive.ObjectID
	Name                  string
	IsAccepted            bool
	EstimatedTotalHarvest float64
	PlantingArea          float64
	Address               string
	IsAvailable           bool
	CreatedAt             primitive.DateTime
	UpdatedAt             primitive.DateTime
	DeletedAt             primitive.DateTime
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByCommodityIDAndName(commodityID primitive.ObjectID, name string) (Domain, error)
	// Update
	// Delete
}

type UseCase interface {
	// Create
	Create(domain *Domain, farmerID primitive.ObjectID) (int, error)
	// Read
	// Update
	// Delete
}
