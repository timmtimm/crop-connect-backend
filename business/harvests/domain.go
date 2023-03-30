package harvests

import (
	"marketplace-backend/dto"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID           primitive.ObjectID
	AccepterID   primitive.ObjectID
	BatchID      primitive.ObjectID
	Date         primitive.DateTime
	Status       string
	TotalHarvest float64
	Condition    string
	Harvest      []dto.ImageAndNote
	RevisionNote string
	CreatedAt    primitive.DateTime
	UpdatedAt    primitive.DateTime
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByBatchID(batchID primitive.ObjectID) (Domain, error)
	// Update
	// Delete
}

type UseCase interface {
	// Create
	SubmitHarvest(domain *Domain, farmerID primitive.ObjectID, images []*multipart.FileHeader, notes []string) (Domain, int, error)
	// Read
	// Update
	// Delete
}
