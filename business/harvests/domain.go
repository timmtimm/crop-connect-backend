package harvests

import (
	"crop_connect/dto"
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

type Query struct {
	Skip      int64
	Limit     int64
	Sort      string
	Order     int
	FarmerID  primitive.ObjectID
	Commodity string
	Batch     string
	Status    string
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByBatchID(batchID primitive.ObjectID) (Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	CountByYear(year int) (float64, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
}

type UseCase interface {
	// Create
	SubmitHarvest(domain *Domain, farmerID primitive.ObjectID, images []*multipart.FileHeader, notes []string) (Domain, int, error)
	// Read
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	GetByBatchID(batchID primitive.ObjectID) (Domain, int, error)
	CountByYear(year int) (float64, int, error)
	// Update
	Validate(domain *Domain, validatorID primitive.ObjectID) (Domain, int, error)
	// Delete
}
