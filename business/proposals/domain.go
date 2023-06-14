package proposals

import (
	"crop_connect/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type Query struct {
	Skip        int64
	Limit       int64
	Sort        string
	Order       int
	CommodityID primitive.ObjectID
	Commodity   string
	FarmerID    primitive.ObjectID
	Name        string
	Status      string
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
	StatisticByYear(year int) ([]dto.StatisticByYear, error)
	CountTotalProposalByFarmer(farmerID primitive.ObjectID) (int, error)
	GetByQuery(query Query) ([]Domain, int, error)
	GetForPerennials(commodityID primitive.ObjectID, farmerID primitive.ObjectID) ([]Domain, error)
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
	StatisticByYear(year int) ([]dto.StatisticByYear, int, error)
	CountTotalProposalByFarmer(farmerID primitive.ObjectID) (int, int, error)
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	GetByIDAndFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (Domain, int, error)
	GetForPerennials(commodityID primitive.ObjectID, farmerID primitive.ObjectID) ([]Domain, int, error)
	// Update
	Update(domain *Domain, farmerID primitive.ObjectID) (int, error)
	UpdateCommodityID(OldCommodityID primitive.ObjectID, NewCommodityID primitive.ObjectID) (int, error)
	ValidateProposal(domain *Domain, adminID primitive.ObjectID) (int, error)
	// Delete
	Delete(id primitive.ObjectID, farmerID primitive.ObjectID) (int, error)
	DeleteByCommodityID(commodityID primitive.ObjectID) (int, error)
}
