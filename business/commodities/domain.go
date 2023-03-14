package commodities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID             primitive.ObjectID `json:"_id"`
	FarmerID       primitive.ObjectID `json:"farmerID"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	Seed           string             `json:"seed"`
	PlantingPeriod int                `json:"plantingPeriod"`
	ImageURLs      []string           `json:"imageURLs"`
	PricePerKg     int                `json:"pricePerKg"`
	IsAvailable    bool               `json:"isAvailable"`
	CreatedAt      primitive.DateTime `json:"createdAt"`
	UpdatedAt      primitive.DateTime `json:"updatedAt"`
	DeletedAt      primitive.DateTime `json:"deletedAt"`
}

type Query struct {
	Skip     int64
	Limit    int64
	Sort     string
	Order    int
	Name     string
	FarmerID []primitive.ObjectID
	MinPrice int
	MaxPrice int
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByIDAndFarmerID(id primitive.ObjectID, farmerID primitive.ObjectID) (Domain, error)
	GetByName(name string) (Domain, error)
	GetByNameAndFarmerID(name string, farmerID primitive.ObjectID) (Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
	Delete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Create(domain *Domain) (int, error)
	// Read
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	GetByID(id primitive.ObjectID) (Domain, int, error)
	// Update
	Update(domain *Domain) (Domain, int, error)
	// Delete
	Delete(id primitive.ObjectID, farmerID primitive.ObjectID) (int, error)
}
