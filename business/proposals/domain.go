package proposals

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	ID                    primitive.ObjectID
	Code                  primitive.ObjectID
	ValidatorID           primitive.ObjectID
	CommodityID           primitive.ObjectID
	RegionID              primitive.ObjectID
	Name                  string
	Description           string
	Status                string
	RejectReason          string
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
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByIDWithoutDeleted(id primitive.ObjectID) (Domain, error)
	GetByCommodityID(commodityID primitive.ObjectID) ([]Domain, error)
	GetByCommodityIDAndAvailability(commodityID primitive.ObjectID, status string) ([]Domain, error)
	GetByCommodityIDAndName(commodityID primitive.ObjectID, name string) (Domain, error)
	GetByIDAccepted(id primitive.ObjectID) (Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	UnsetRejectReason(id primitive.ObjectID) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain, farmerID primitive.ObjectID) (int, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, int, error)
	GetByCommodityID(commodityID primitive.ObjectID) ([]Domain, int, error)
	GetByIDWithoutDeleted(id primitive.ObjectID) (Domain, int, error)
	GetByIDAccepted(id primitive.ObjectID) (Domain, int, error)
	// Update
	Update(domain *Domain, farmerID primitive.ObjectID) (int, error)
	UpdateCommodityID(OldCommodityID primitive.ObjectID, NewCommodityID primitive.ObjectID) (int, error)
	ValidateProposal(domain *Domain, adminID primitive.ObjectID) (int, error)
	// Delete
	Delete(id primitive.ObjectID, farmerID primitive.ObjectID) (int, error)
	DeleteByCommodityID(commodityID primitive.ObjectID) (int, error)
}
