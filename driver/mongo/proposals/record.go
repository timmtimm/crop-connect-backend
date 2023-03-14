package proposals

import (
	"marketplace-backend/business/proposals"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID                    primitive.ObjectID `bson:"_id"`
	AccepterID            primitive.ObjectID `bson:"accepterID,omitempty"`
	CommodityID           primitive.ObjectID `bson:"commodityID"`
	Name                  string             `bson:"name"`
	Description           string             `bson:"description"`
	IsAccepted            bool               `bson:"isAccepted"`
	EstimatedTotalHarvest float64            `bson:"estimatedTotalHarvest"`
	PlantingArea          float64            `bson:"plantingArea"`
	Address               string             `bson:"address"`
	IsAvailable           bool               `bson:"isAvailable,omitempty"`
	CreatedAt             primitive.DateTime `bson:"createdAt"`
	UpdatedAt             primitive.DateTime `bson:"updatedAt,omitempty"`
	DeletedAt             primitive.DateTime `bson:"deletedAt,omitempty"`
}

func FromDomain(domain *proposals.Domain) *Model {
	return &Model{
		ID:                    domain.ID,
		AccepterID:            domain.AccepterID,
		CommodityID:           domain.CommodityID,
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
	}
}

func (model *Model) ToDomain() proposals.Domain {
	return proposals.Domain{
		ID:                    model.ID,
		AccepterID:            model.AccepterID,
		CommodityID:           model.CommodityID,
		Name:                  model.Name,
		Description:           model.Description,
		IsAccepted:            model.IsAccepted,
		EstimatedTotalHarvest: model.EstimatedTotalHarvest,
		PlantingArea:          model.PlantingArea,
		Address:               model.Address,
		IsAvailable:           model.IsAvailable,
		CreatedAt:             model.CreatedAt,
		UpdatedAt:             model.UpdatedAt,
		DeletedAt:             model.DeletedAt,
	}
}

func ToDomainArray(model []Model) []proposals.Domain {
	var result []proposals.Domain
	for _, proposal := range model {
		result = append(result, proposal.ToDomain())
	}
	return result
}
