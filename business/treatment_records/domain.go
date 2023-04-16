package treatment_records

import (
	"marketplace-backend/dto"
	"mime/multipart"

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
	Treatment    []dto.ImageAndNote
	RevisionNote string
	WarningNote  string
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
	Number    int
	Status    string
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetNewestByBatchID(batchID primitive.ObjectID) (Domain, error)
	CountByBatchID(batchID primitive.ObjectID) (int, error)
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByBatchID(batchID primitive.ObjectID) ([]Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
}

type UseCase interface {
	// Create
	RequestToFarmer(domain *Domain) (Domain, int, error)
	// Read
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	GetByBatchID(batchID primitive.ObjectID) ([]Domain, int, error)
	// Update
	FillTreatmentRecord(domain *Domain, farmerID primitive.ObjectID, images []*multipart.FileHeader, notes []string) (Domain, int, error)
	Validate(domain *Domain, validatorID primitive.ObjectID) (Domain, int, error)
	UpdateNotes(domain *Domain) (Domain, int, error)
	// Delete
}
