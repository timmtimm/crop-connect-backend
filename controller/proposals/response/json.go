package response

import (
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/proposals"
	"marketplace-backend/business/users"
	commodityResponse "marketplace-backend/controller/commodities/response"
	userReponse "marketplace-backend/controller/users/response"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	ID                    primitive.ObjectID          `json:"_id"`
	Accepter              userReponse.User            `json:"accepter"`
	Commodity             commodityResponse.Commodity `json:"commodity"`
	Name                  string                      `json:"name"`
	Description           string                      `json:"description"`
	IsAccepted            bool                        `json:"isAccepted"`
	EstimatedTotalHarvest float64                     `json:"estimatedTotalHarvest"`
	PlantingArea          float64                     `json:"plantingArea"`
	Address               string                      `json:"address"`
	IsAvailable           bool                        `json:"isAvailable"`
	CreatedAt             primitive.DateTime          `json:"createdAt"`
	UpdatedAt             primitive.DateTime          `json:"updatedAt"`
	DeletedAt             primitive.DateTime          `json:"deletedAt"`
}

func FromDomainToAdmin(domain *proposals.Domain, userUC users.UseCase, commodityUC commodities.UseCase) (Admin, int, error) {
	accepter, statusCode, err := userUC.GetByID(domain.AccepterID)
	if err != nil {
		return Admin{}, statusCode, err
	}

	commodity, statusCode, err := commodityUC.GetByID(domain.CommodityID)
	if err != nil {
		return Admin{}, statusCode, err
	}

	commodityResponse, statusCode, err := commodityResponse.FromDomain(commodity, userUC)
	if err != nil {
		return Admin{}, statusCode, err
	}

	return Admin{
		ID:                    domain.ID,
		Accepter:              userReponse.FromDomain(accepter),
		Commodity:             commodityResponse,
		Name:                  domain.Name,
		Description:           domain.Description,
		IsAccepted:            domain.IsAccepted,
		EstimatedTotalHarvest: domain.EstimatedTotalHarvest,
		PlantingArea:          domain.PlantingArea,
		Address:               domain.Address,
		IsAvailable:           domain.IsAvailable,
		CreatedAt:             domain.CreatedAt,
		UpdatedAt:             domain.UpdatedAt,
		DeletedAt:             domain.DeletedAt,
	}, http.StatusOK, nil
}
