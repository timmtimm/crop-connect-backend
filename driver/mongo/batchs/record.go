package batchs

import (
	"crop_connect/business/batchs"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID                   primitive.ObjectID `bson:"_id"`
	ProposalID           primitive.ObjectID `bson:"proposalID"`
	Name                 string             `bson:"name"`
	EstimatedHarvestDate primitive.DateTime `bson:"estimatedHarvestDate"`
	Status               string             `bson:"status"`
	CancelReason         string             `bson:"cancelReason,omitempty"`
	IsAvailable          bool               `bson:"isAvailable"`
	CreatedAt            primitive.DateTime `bson:"createdAt"`
	UpdatedAt            primitive.DateTime `bson:"updatedAt,omitempty"`
}

func FromDomain(domain *batchs.Domain) *Model {
	return &Model{
		ID:                   domain.ID,
		ProposalID:           domain.ProposalID,
		Name:                 domain.Name,
		EstimatedHarvestDate: domain.EstimatedHarvestDate,
		Status:               domain.Status,
		CancelReason:         domain.CancelReason,
		IsAvailable:          domain.IsAvailable,
		CreatedAt:            domain.CreatedAt,
		UpdatedAt:            domain.UpdatedAt,
	}
}

func (model *Model) ToDomain() batchs.Domain {
	return batchs.Domain{
		ID:                   model.ID,
		ProposalID:           model.ProposalID,
		Name:                 model.Name,
		EstimatedHarvestDate: model.EstimatedHarvestDate,
		Status:               model.Status,
		CancelReason:         model.CancelReason,
		IsAvailable:          model.IsAvailable,
		CreatedAt:            model.CreatedAt,
		UpdatedAt:            model.UpdatedAt,
	}
}

func ToDomainArray(model []Model) []batchs.Domain {
	var domain []batchs.Domain
	for _, v := range model {
		domain = append(domain, v.ToDomain())
	}
	return domain
}
