package response

import (
	"crop_connect/business/commodities"
	"crop_connect/business/regions"
	"crop_connect/business/users"
	userReponse "crop_connect/controller/users/response"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Commodity struct {
	ID             primitive.ObjectID `json:"_id"`
	Code           primitive.ObjectID `json:"code"`
	Farmer         userReponse.User   `json:"farmer"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	Seed           string             `json:"seed"`
	PlantingPeriod int                `json:"plantingPeriod"`
	ImageURLs      []string           `json:"imageURLs"`
	PricePerKg     int                `json:"pricePerKg"`
	IsAvailable    bool               `json:"isAvailable"`
	CreatedAt      primitive.DateTime `json:"createdAt"`
	UpdatedAt      primitive.DateTime `json:"updatedAt,omitempty"`
	DeletedAt      primitive.DateTime `json:"deletedAt,omitempty"`
}

func FromDomain(domain commodities.Domain, userUC users.UseCase, regionUC regions.UseCase) (Commodity, int, error) {
	farmer, statusCode, err := userUC.GetByID(domain.FarmerID)
	if err != nil {
		return Commodity{}, statusCode, err
	}

	farmerResponse, statusCode, err := userReponse.FromDomain(farmer, regionUC)
	if err != nil {
		return Commodity{}, statusCode, err
	}

	return Commodity{
		ID:             domain.ID,
		Code:           domain.Code,
		Farmer:         farmerResponse,
		Name:           domain.Name,
		Description:    domain.Description,
		Seed:           domain.Seed,
		PlantingPeriod: domain.PlantingPeriod,
		ImageURLs:      domain.ImageURLs,
		PricePerKg:     domain.PricePerKg,
		IsAvailable:    domain.IsAvailable,
		CreatedAt:      domain.CreatedAt,
		UpdatedAt:      domain.UpdatedAt,
		DeletedAt:      domain.DeletedAt,
	}, http.StatusOK, nil
}

func FromDomainArray(domain []commodities.Domain, userUC users.UseCase, regionUC regions.UseCase) ([]Commodity, int, error) {
	var response []Commodity
	for _, value := range domain {
		commodity, statusCode, err := FromDomain(value, userUC, regionUC)
		if err != nil {
			return []Commodity{}, statusCode, err
		}

		response = append(response, commodity)
	}

	return response, http.StatusOK, nil
}
