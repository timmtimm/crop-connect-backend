package treatment_records

import (
	treatmentRecord "marketplace-backend/business/treatment_records"
	"marketplace-backend/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID           primitive.ObjectID `bson:"_id"`
	RequesterID  primitive.ObjectID `bson:"requesterID"`
	AccepterID   primitive.ObjectID `bson:"accepterID,omitempty"`
	BatchID      primitive.ObjectID `bson:"batchID"`
	Number       int                `bson:"number"`
	Date         primitive.DateTime `bson:"date"`
	Status       string             `bson:"status"`
	Description  string             `bson:"description"`
	Treatment    []dto.ImageAndNote `bson:"treatment,omitempty"`
	RevisionNote string             `bson:"revisionNote,omitempty"`
	WarningNote  string             `bson:"warningNote,omitempty"`
	CreatedAt    primitive.DateTime `bson:"createdAt"`
	UpdatedAt    primitive.DateTime `bson:"updatedAt,omitempty"`
}

func FromDomain(domain *treatmentRecord.Domain) *Model {
	return &Model{
		ID:           domain.ID,
		RequesterID:  domain.RequesterID,
		AccepterID:   domain.AccepterID,
		BatchID:      domain.BatchID,
		Number:       domain.Number,
		Date:         domain.Date,
		Status:       domain.Status,
		Description:  domain.Description,
		Treatment:    domain.Treatment,
		RevisionNote: domain.RevisionNote,
		WarningNote:  domain.WarningNote,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
	}
}

func (model *Model) ToDomain() treatmentRecord.Domain {
	return treatmentRecord.Domain{
		ID:           model.ID,
		RequesterID:  model.RequesterID,
		AccepterID:   model.AccepterID,
		BatchID:      model.BatchID,
		Number:       model.Number,
		Date:         model.Date,
		Status:       model.Status,
		Description:  model.Description,
		Treatment:    model.Treatment,
		RevisionNote: model.RevisionNote,
		WarningNote:  model.WarningNote,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func ToDomainArray(models []Model) []treatmentRecord.Domain {
	var domains []treatmentRecord.Domain
	for _, model := range models {
		domains = append(domains, model.ToDomain())
	}
	return domains
}
