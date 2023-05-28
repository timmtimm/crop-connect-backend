package response

import (
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/users"
	commodityResponse "crop_connect/controller/commodities/response"
	regionResponse "crop_connect/controller/regions/response"
	userReponse "crop_connect/controller/users/response"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	ID                    primitive.ObjectID          `json:"_id"`
	Validator             userReponse.User            `json:"validator"`
	Commodity             commodityResponse.Commodity `json:"commodity"`
	Name                  string                      `json:"name"`
	Description           string                      `json:"description"`
	Status                string                      `json:"status"`
	EstimatedTotalHarvest float64                     `json:"estimatedTotalHarvest"`
	PlantingArea          float64                     `json:"plantingArea"`
	Address               string                      `json:"address"`
	IsAvailable           bool                        `json:"isAvailable"`
	CreatedAt             primitive.DateTime          `json:"createdAt"`
	UpdatedAt             primitive.DateTime          `json:"updatedAt,omitempty"`
	DeletedAt             primitive.DateTime          `json:"deletedAt,omitempty"`
}

func FromDomainToAdmin(domain *proposals.Domain, userUC users.UseCase, commodityUC commodities.UseCase, regionUC regions.UseCase) (Admin, int, error) {
	validator, statusCode, err := userUC.GetByID(domain.ValidatorID)
	if err != nil {
		return Admin{}, statusCode, err
	}

	commodity, statusCode, err := commodityUC.GetByID(domain.CommodityID)
	if err != nil {
		return Admin{}, statusCode, err
	}

	commodityResponse, statusCode, err := commodityResponse.FromDomain(commodity, userUC, regionUC)
	if err != nil {
		return Admin{}, statusCode, err
	}

	validatorResponse, statusCode, err := userReponse.FromDomain(validator, regionUC)
	if err != nil {
		return Admin{}, statusCode, err
	}

	return Admin{
		ID:                    domain.ID,
		Validator:             validatorResponse,
		Commodity:             commodityResponse,
		Name:                  domain.Name,
		Description:           domain.Description,
		Status:                domain.Status,
		EstimatedTotalHarvest: domain.EstimatedTotalHarvest,
		PlantingArea:          domain.PlantingArea,
		Address:               domain.Address,
		IsAvailable:           domain.IsAvailable,
		CreatedAt:             domain.CreatedAt,
		UpdatedAt:             domain.UpdatedAt,
		DeletedAt:             domain.DeletedAt,
	}, http.StatusOK, nil
}

type Buyer struct {
	ID                    primitive.ObjectID `json:"_id"`
	Name                  string             `json:"name"`
	Description           string             `json:"description"`
	EstimatedTotalHarvest float64            `json:"estimatedTotalHarvest"`
	PlantingArea          float64            `json:"plantingArea"`
	Address               string             `json:"address"`
	IsAvailable           bool               `json:"isAvailable"`
}

func FromDomainToBuyer(domain *proposals.Domain) Buyer {
	return Buyer{
		ID:                    domain.ID,
		Name:                  domain.Name,
		Description:           domain.Description,
		EstimatedTotalHarvest: domain.EstimatedTotalHarvest,
		PlantingArea:          domain.PlantingArea,
		Address:               domain.Address,
		IsAvailable:           domain.IsAvailable,
	}
}

func FromDomainArrayToBuyer(domain []proposals.Domain) []Buyer {
	var response []Buyer
	for _, value := range domain {
		response = append(response, FromDomainToBuyer(&value))
	}

	return response
}

type ProposalWithCommodity struct {
	ID                    primitive.ObjectID          `json:"_id"`
	Commodity             commodityResponse.Commodity `json:"commodity"`
	Region                regionResponse.Response     `json:"region"`
	Name                  string                      `json:"name"`
	Description           string                      `json:"description"`
	EstimatedTotalHarvest float64                     `json:"estimatedTotalHarvest"`
	PlantingArea          float64                     `json:"plantingArea"`
	Address               string                      `json:"address"`
	IsAvailable           bool                        `json:"isAvailable"`
}

func FromDomainToProposalWithCommodity(domain *proposals.Domain, userUC users.UseCase, commodityUC commodities.UseCase, regionUC regions.UseCase) (ProposalWithCommodity, int, error) {
	commodity, statusCode, err := commodityUC.GetByIDWithoutDeleted(domain.CommodityID)
	if err != nil {
		return ProposalWithCommodity{}, statusCode, err
	}

	commodityResponse, statusCode, err := commodityResponse.FromDomain(commodity, userUC, regionUC)
	if err != nil {
		return ProposalWithCommodity{}, statusCode, err
	}

	region, statusCode, err := regionUC.GetByID(domain.RegionID)
	if err != nil {
		return ProposalWithCommodity{}, statusCode, err
	}

	return ProposalWithCommodity{
		ID:                    domain.ID,
		Commodity:             commodityResponse,
		Region:                regionResponse.FromDomain(&region),
		Name:                  domain.Name,
		Description:           domain.Description,
		EstimatedTotalHarvest: domain.EstimatedTotalHarvest,
		PlantingArea:          domain.PlantingArea,
		Address:               domain.Address,
		IsAvailable:           domain.IsAvailable,
	}, http.StatusOK, nil
}

func FromDomainArrayToProposalWithCommodity(domain []proposals.Domain, userUC users.UseCase, commodityUC commodities.UseCase, regionUC regions.UseCase) ([]ProposalWithCommodity, int, error) {
	var response []ProposalWithCommodity
	for _, value := range domain {
		proposal, statusCode, err := FromDomainToProposalWithCommodity(&value, userUC, commodityUC, regionUC)
		if err != nil {
			return []ProposalWithCommodity{}, statusCode, err
		}

		response = append(response, proposal)
	}

	return response, http.StatusOK, nil
}
