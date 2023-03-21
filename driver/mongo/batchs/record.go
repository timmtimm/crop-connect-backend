package batchs

import (
	"marketplace-backend/business/batchs"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID                   primitive.ObjectID `bson:"_id"`
	TransactionID        primitive.ObjectID `bson:"transactionID"`
	Name                 string             `bson:"name"`
	EstimatedHarvestDate primitive.DateTime `bson:"estimatedHarvestDate"`
	Status               string             `bson:"status"`
	CancelReason         string             `bson:"cancelReason,omitempty"`
	CreatedAt            primitive.DateTime `bson:"createdAt"`
	UpdatedAt            primitive.DateTime `bson:"updatedAt,omitempty"`
}

func FromDomain(domain *batchs.Domain) *Model {
	return &Model{
		ID:                   domain.ID,
		TransactionID:        domain.TransactionID,
		Name:                 domain.Name,
		EstimatedHarvestDate: domain.EstimatedHarvestDate,
		Status:               domain.Status,
		CancelReason:         domain.CancelReason,
		CreatedAt:            domain.CreatedAt,
		UpdatedAt:            domain.UpdatedAt,
	}
}

func (model *Model) ToDomain() *batchs.Domain {
	return &batchs.Domain{
		ID:                   model.ID,
		TransactionID:        model.TransactionID,
		Name:                 model.Name,
		EstimatedHarvestDate: model.EstimatedHarvestDate,
		Status:               model.Status,
		CancelReason:         model.CancelReason,
		CreatedAt:            model.CreatedAt,
		UpdatedAt:            model.UpdatedAt,
	}
}

func ToDomainArray(model []Model) []batchs.Domain {
	var domain []batchs.Domain
	for _, v := range model {
		domain = append(domain, *v.ToDomain())
	}
	return domain
}
