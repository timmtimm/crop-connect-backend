package harvests

import (
	"crop_connect/dto"
	"crop_connect/helper"
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
	Skip        int64
	Limit       int64
	Sort        string
	Order       int
	FarmerID    primitive.ObjectID
	CommodityID primitive.ObjectID
	Commodity   string
	BatchID     primitive.ObjectID
	Batch       string
	Status      string
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByBatchIDAndStatus(batchID primitive.ObjectID, status string) (Domain, error)
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
	GetByBatchIDAndStatus(batchID primitive.ObjectID, status string) (Domain, int, error)
	CountByYear(year int) (float64, int, error)
	GetByID(id primitive.ObjectID) (Domain, int, error)
	// Update
	UpdateHarvest(domain *Domain, farmerID primitive.ObjectID, updateImages []*helper.UpdateImage, notes []string) (Domain, int, error)
	Validate(domain *Domain, validatorID primitive.ObjectID) (Domain, int, error)
	// Delete
}
