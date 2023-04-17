package harvests

import (
	"crop_connect/business/harvests"
	"crop_connect/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID           primitive.ObjectID `bson:"_id"`
	AccepterID   primitive.ObjectID `bson:"accepterID,omitempty"`
	BatchID      primitive.ObjectID `bson:"batchID"`
	Date         primitive.DateTime `bson:"date"`
	Status       string             `bson:"status"`
	TotalHarvest float64            `bson:"totalHarvest"`
	Condition    string             `bson:"condition"`
	Harvest      []dto.ImageAndNote `bson:"harvest"`
	RevisionNote string             `bson:"revisionNote,omitempty"`
	CreatedAt    primitive.DateTime `bson:"createdAt"`
	UpdatedAt    primitive.DateTime `bson:"updatedAt,omitempty"`
}

func FromDomain(domain *harvests.Domain) *Model {
	return &Model{
		ID:           domain.ID,
		AccepterID:   domain.AccepterID,
		BatchID:      domain.BatchID,
		Date:         domain.Date,
		Status:       domain.Status,
		TotalHarvest: domain.TotalHarvest,
		Condition:    domain.Condition,
		Harvest:      domain.Harvest,
		RevisionNote: domain.RevisionNote,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
	}
}

func (model *Model) ToDomain() harvests.Domain {
	return harvests.Domain{
		ID:           model.ID,
		AccepterID:   model.AccepterID,
		BatchID:      model.BatchID,
		Date:         model.Date,
		Status:       model.Status,
		TotalHarvest: model.TotalHarvest,
		Condition:    model.Condition,
		Harvest:      model.Harvest,
		RevisionNote: model.RevisionNote,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func ToDomainArray(models []Model) []harvests.Domain {
	var domains []harvests.Domain
	for _, model := range models {
		domains = append(domains, model.ToDomain())
	}
	return domains
}
