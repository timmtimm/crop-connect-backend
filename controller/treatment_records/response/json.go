package response

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/transactions"
	treatmentRecord "crop_connect/business/treatment_records"
	"crop_connect/business/users"
	batchResponse "crop_connect/controller/batchs/response"
	proposalResponse "crop_connect/controller/proposals/response"
	userResponse "crop_connect/controller/users/response"
	"crop_connect/dto"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TreatmentRecord struct {
	ID           primitive.ObjectID                     `json:"_id"`
	Requester    userResponse.User                      `json:"requester"`
	Accepter     interface{}                            `json:"accepter,omitempty"`
	Proposal     proposalResponse.ProposalWithCommodity `json:"proposal"`
	Batch        batchResponse.BatchWithoutProposal     `json:"batch"`
	Number       int                                    `json:"number"`
	Date         primitive.DateTime                     `json:"date"`
	Status       string                                 `json:"status"`
	Description  string                                 `json:"description"`
	Treatment    []dto.ImageAndNote                     `json:"treatment,omitempty"`
	RevisionNote string                                 `json:"revisionNote,omitempty"`
	WarningNote  string                                 `json:"warningNote,omitempty"`
	CreatedAt    primitive.DateTime                     `json:"createdAt"`
	UpdatedAt    primitive.DateTime                     `json:"updatedAt,omitempty"`
}

func FromDomain(domain treatmentRecord.Domain, batchUC batchs.UseCase, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) (*TreatmentRecord, int, error) {
	batch, statusCode, err := batchUC.GetByID(domain.BatchID)
	if err != nil {
		return nil, statusCode, err
	}

	proposal, statusCode, err := proposalUC.GetByIDWithoutDeleted(batch.ProposalID)
	if err != nil {
		return nil, statusCode, err
	}

	proposalResponse, statusCode, err := proposalResponse.FromDomainToProposalWithCommodity(&proposal, userUC, commodityUC, regionUC)
	if err != nil {
		return nil, statusCode, err
	}

	requester, statusCode, err := userUC.GetByID(domain.RequesterID)
	if err != nil {
		return nil, statusCode, err
	}

	requesterResponse, statusCode, err := userResponse.FromDomain(requester, regionUC)
	if err != nil {
		return nil, statusCode, err
	}

	response := TreatmentRecord{
		ID:           domain.ID,
		Requester:    requesterResponse,
		Proposal:     proposalResponse,
		Batch:        batchResponse.FromDomainWithoutProposal(&batch),
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

	var accepter users.Domain
	if !domain.AccepterID.IsZero() {
		accepter, statusCode, err = userUC.GetByID(domain.AccepterID)
		if err != nil {
			return nil, statusCode, err
		}

		accepterResponse, statusCode, err := userResponse.FromDomain(accepter, regionUC)
		if err != nil {
			return nil, statusCode, err
		}

		response.Accepter = accepterResponse
	}

	return &response, http.StatusOK, nil
}

func FromDomainArray(domain []treatmentRecord.Domain, batchUC batchs.UseCase, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) ([]TreatmentRecord, int, error) {
	var treatmentRecords []TreatmentRecord
	for _, v := range domain {
		treatmentRecord, statusCode, err := FromDomain(v, batchUC, transactionUC, proposalUC, commodityUC, userUC, regionUC)
		if err != nil {
			return treatmentRecords, statusCode, err
		}

		treatmentRecords = append(treatmentRecords, *treatmentRecord)
	}

	return treatmentRecords, http.StatusOK, nil
}
