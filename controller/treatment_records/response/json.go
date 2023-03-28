package response

import (
	"marketplace-backend/business/batchs"
	treatmentRecord "marketplace-backend/business/treatment_records"
	"marketplace-backend/business/users"
	batchResponse "marketplace-backend/controller/batchs/response"
	userResponse "marketplace-backend/controller/users/response"
	"marketplace-backend/dto"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TreatmentRecord struct {
	ID           primitive.ObjectID                    `json:"id"`
	Requester    userResponse.User                     `json:"requester"`
	Accepter     interface{}                           `json:"accepter,omitempty"`
	Batch        batchResponse.BatchWithoutTransaction `json:"batch"`
	Number       int                                   `json:"number"`
	Date         primitive.DateTime                    `json:"date"`
	Status       string                                `json:"status"`
	Description  string                                `json:"description"`
	Treatment    []dto.Treatment                       `json:"treatment,omitempty"`
	RevisionNote string                                `json:"revisionNote,omitempty"`
	WarningNote  string                                `json:"warningNote,omitempty"`
	CreatedAt    primitive.DateTime                    `json:"createdAt"`
	UpdatedAt    primitive.DateTime                    `json:"updatedAt,omitempty"`
}

func FromDomain(domain treatmentRecord.Domain, batchUC batchs.UseCase, userUC users.UseCase) (*TreatmentRecord, int, error) {
	requester, statusCode, err := userUC.GetByID(domain.RequesterID)
	if err != nil {
		return nil, statusCode, err
	}

	batch, statusCode, err := batchUC.GetByID(domain.BatchID)
	if err != nil {
		return nil, statusCode, err
	}

	var accepter users.Domain
	if !domain.AccepterID.IsZero() {
		accepter, statusCode, err = userUC.GetByID(domain.AccepterID)
		if err != nil {
			return nil, statusCode, err
		}

		return &TreatmentRecord{
			ID:           domain.ID,
			Requester:    userResponse.FromDomain(requester),
			Accepter:     userResponse.FromDomain(accepter),
			Batch:        batchResponse.FromDomainWithoutTransaction(&batch),
			Number:       domain.Number,
			Date:         domain.Date,
			Status:       domain.Status,
			Description:  domain.Description,
			Treatment:    domain.Treatment,
			RevisionNote: domain.RevisionNote,
			WarningNote:  domain.WarningNote,
			CreatedAt:    domain.CreatedAt,
			UpdatedAt:    domain.UpdatedAt,
		}, http.StatusOK, nil
	}

	return &TreatmentRecord{
		ID:           domain.ID,
		Requester:    userResponse.FromDomain(requester),
		Batch:        batchResponse.FromDomainWithoutTransaction(&batch),
		Number:       domain.Number,
		Date:         domain.Date,
		Status:       domain.Status,
		Description:  domain.Description,
		Treatment:    domain.Treatment,
		RevisionNote: domain.RevisionNote,
		WarningNote:  domain.WarningNote,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
	}, http.StatusOK, nil
}

func FromDomainArray(domain []treatmentRecord.Domain, batchUC batchs.UseCase, userUC users.UseCase) ([]TreatmentRecord, int, error) {
	var treatmentRecords []TreatmentRecord
	for _, v := range domain {
		treatmentRecord, statusCode, err := FromDomain(v, batchUC, userUC)
		if err != nil {
			return nil, statusCode, err
		}
		treatmentRecords = append(treatmentRecords, *treatmentRecord)
	}

	return treatmentRecords, http.StatusOK, nil
}
