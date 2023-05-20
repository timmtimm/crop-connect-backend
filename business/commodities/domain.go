package commodities

import (
	"crop_connect/helper"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID             primitive.ObjectID
	Code           primitive.ObjectID
	FarmerID       primitive.ObjectID
	Name           string
	Description    string
	Seed           string
	PlantingPeriod int
	ImageURLs      []string
	PricePerKg     int
	IsAvailable    bool
	CreatedAt      primitive.DateTime
	UpdatedAt      primitive.DateTime
	DeletedAt      primitive.DateTime
}

type Query struct {
	Skip     int64
	Limit    int64
	Sort     string
	Order    int
	Name     string
	Farmer   string
	FarmerID primitive.ObjectID
	MinPrice int
	MaxPrice int
	Province string
	Regency  string
	District string
	RegionID primitive.ObjectID
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByIDWithoutDeleted(id primitive.ObjectID) (Domain, error)
	GetByIDAndFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (Domain, error)
	GetByName(name string) (Domain, error)
	GetByNameAndFarmerID(name string, farmerID primitive.ObjectID) (Domain, error)
	GetByFarmerID(farmerID primitive.ObjectID) ([]Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	CountTotalCommodity(year int) (int, error)
	CountTotalCommodityByFarmer(farmerID primitive.ObjectID) (int, error)
	GetByCode(code primitive.ObjectID) (Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain, images []*multipart.FileHeader) (int, error)
	// Read
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	GetByID(id primitive.ObjectID) (Domain, int, error)
	GetByIDWithoutDeleted(id primitive.ObjectID) (Domain, int, error)
	GetByFarmerID(farmerID primitive.ObjectID) ([]Domain, int, error)
	GetByIDAndFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (Domain, int, error)
	CountTotalCommodityByFarmer(farmerID primitive.ObjectID) (int, int, error)
	CountTotalCommodity(year int) (int, int, error)
	// Update
	Update(domain *Domain, updateImage []*helper.UpdateImage) (Domain, int, error)
	// Delete
	Delete(id primitive.ObjectID, farmerID primitive.ObjectID) (int, error)
}
