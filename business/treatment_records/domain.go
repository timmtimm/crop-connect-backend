package treatment_records

import (
	"marketplace-backend/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID           primitive.ObjectID
	RequesterID  primitive.ObjectID
	AccepterID   primitive.ObjectID
	BatchID      primitive.ObjectID
	Number       int
	Date         primitive.DateTime
	Status       string
	Description  string
	Treatment    []dto.Treatment
	RevisionNote string
	WarningNote  string
	CreatedAt    primitive.DateTime
	UpdatedAt    primitive.DateTime
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetNewestByBatchID(batchID primitive.ObjectID) (Domain, error)
	CountByBatchID(batchID primitive.ObjectID) (int, error)
	// Update
	// Delete
}

type UseCase interface {
	// Create
	RequestToFarmer(domain *Domain) (Domain, int, error)
	// Read
	// Update
	// Delete
}
