package commodities

import (
	"crop_connect/business/commodities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID             primitive.ObjectID `bson:"_id"`
	Code           primitive.ObjectID `bson:"code"`
	FarmerID       primitive.ObjectID `bson:"farmerID"`
	Name           string             `bson:"name"`
	Description    string             `bson:"description"`
	Seed           string             `bson:"seed"`
	PlantingPeriod int                `bson:"plantingPeriod"`
	ImageURLs      []string           `bson:"imageURLs"`
	PricePerKg     int                `bson:"pricePerKg"`
	IsPerennials   bool               `bson:"isPerennials"`
	IsAvailable    bool               `bson:"isAvailable"`
	CreatedAt      primitive.DateTime `bson:"createdAt"`
	UpdatedAt      primitive.DateTime `bson:"updatedAt,omitempty"`
	DeletedAt      primitive.DateTime `bson:"deletedAt,omitempty"`
}

func FromDomain(domain *commodities.Domain) *Model {
	return &Model{
		ID:             domain.ID,
		Code:           domain.Code,
		FarmerID:       domain.FarmerID,
		Name:           domain.Name,
		Description:    domain.Description,
		Seed:           domain.Seed,
		PlantingPeriod: domain.PlantingPeriod,
		ImageURLs:      domain.ImageURLs,
		PricePerKg:     domain.PricePerKg,
		IsPerennials:   domain.IsPerennials,
		IsAvailable:    domain.IsAvailable,
		CreatedAt:      domain.CreatedAt,
		UpdatedAt:      domain.UpdatedAt,
		DeletedAt:      domain.DeletedAt,
	}
}

func (model *Model) ToDomain() commodities.Domain {
	return commodities.Domain{
		ID:             model.ID,
		Code:           model.Code,
		FarmerID:       model.FarmerID,
		Name:           model.Name,
		Description:    model.Description,
		Seed:           model.Seed,
		PlantingPeriod: model.PlantingPeriod,
		ImageURLs:      model.ImageURLs,
		PricePerKg:     model.PricePerKg,
		IsPerennials:   model.IsPerennials,
		IsAvailable:    model.IsAvailable,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		DeletedAt:      model.DeletedAt,
	}
}

func ToDomainArray(models []Model) []commodities.Domain {
	var domains []commodities.Domain
	for _, model := range models {
		domains = append(domains, model.ToDomain())
	}
	return domains
}
