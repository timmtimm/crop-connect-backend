package response

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/harvests"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/transactions"
	"crop_connect/business/users"
	batchResponse "crop_connect/controller/batchs/response"
	proposalResponse "crop_connect/controller/proposals/response"
	userResponse "crop_connect/controller/users/response"
	"crop_connect/dto"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Harvest struct {
	ID           string                                 `json:"_id"`
	Accepter     interface{}                            `json:"accepter,omitempty"`
	Proposal     proposalResponse.ProposalWithCommodity `json:"proposal"`
	Batch        batchResponse.BatchWithoutProposal     `json:"batch"`
	Date         primitive.DateTime                     `json:"date"`
	Status       string                                 `json:"status"`
	TotalHarvest float64                                `json:"totalHarvest"`
	Condition    string                                 `json:"condition"`
	Harvest      []dto.ImageAndNote                     `json:"harvest"`
	RevisionNote string                                 `json:"revisionNote"`
	CreatedAt    primitive.DateTime                     `json:"createdAt"`
	UpdatedAt    primitive.DateTime                     `json:"updatedAt,omitempty"`
}

func FromDomain(domain *harvests.Domain, batchUC batchs.UseCase, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) (Harvest, int, error) {
	harvestResponse := Harvest{
		ID:           domain.ID.Hex(),
		Date:         domain.Date,
		Status:       domain.Status,
		TotalHarvest: domain.TotalHarvest,
		Condition:    domain.Condition,
		Harvest:      domain.Harvest,
		RevisionNote: domain.RevisionNote,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
	}

	batch, statusCode, err := batchUC.GetByID(domain.BatchID)
	if err != nil {
		return Harvest{}, statusCode, err
	}

	harvestResponse.Batch = batchResponse.FromDomainWithoutProposal(&batch)

	proposal, statusCode, err := proposalUC.GetByID(batch.ProposalID)
	if err != nil {
		return Harvest{}, statusCode, err
	}

	proposalResponse, statusCode, err := proposalResponse.FromDomainToProposalWithCommodity(&proposal, userUC, commodityUC, regionUC)
	if err != nil {
		return Harvest{}, statusCode, err
	}

	harvestResponse.Proposal = proposalResponse

	if domain.AccepterID != primitive.NilObjectID {
		accepter, statusCode, err := userUC.GetByID(domain.AccepterID)
		if err != nil {
			return Harvest{}, statusCode, err
		}

		accepterResponse, statusCode, err := userResponse.FromDomain(accepter, regionUC)
		if err != nil {
			return Harvest{}, statusCode, err
		}

		harvestResponse.Accepter = accepterResponse
	}

	return harvestResponse, http.StatusOK, nil
}

func FromDomainArrayToResponse(domain []harvests.Domain, batchUC batchs.UseCase, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) ([]Harvest, int, error) {
	response := []Harvest{}

	for _, v := range domain {
		harvestResponse, statusCode, err := FromDomain(&v, batchUC, transactionUC, proposalUC, commodityUC, userUC, regionUC)
		if err != nil {
			return []Harvest{}, statusCode, err
		}

		response = append(response, harvestResponse)
	}

	return response, http.StatusOK, nil
}
