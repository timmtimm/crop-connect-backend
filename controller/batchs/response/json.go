package response

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/users"
	"errors"
	"net/http"

	proposalResponse "crop_connect/controller/proposals/response"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Batch struct {
	ID                   primitive.ObjectID                     `json:"_id"`
	Proposal             proposalResponse.ProposalWithCommodity `json:"proposal"`
	Name                 string                                 `json:"name"`
	EstimatedHarvestDate primitive.DateTime                     `json:"estimatedHarvestDate"`
	Status               string                                 `json:"status"`
	CancelReason         string                                 `json:"cancelReason,omitempty"`
	IsAvailable          bool                                   `json:"isAvailable"`
	CreatedAt            primitive.DateTime                     `json:"createdAt"`
	UpdatedAt            primitive.DateTime                     `json:"updatedAt,omitempty"`
}

func FromDomain(domain batchs.Domain, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) (Batch, int, error) {
	proposal, statusCode, err := proposalUC.GetByID(domain.ProposalID)
	if err != nil {
		return Batch{}, statusCode, errors.New("gagal mendapatkan transaksi")
	}

	responseForProposal, statusCode, err := proposalResponse.FromDomainToProposalWithCommodity(&proposal, userUC, commodityUC, regionUC)
	if err != nil {
		return Batch{}, statusCode, errors.New("gagal mendapatkan transaksi")
	}

	return Batch{
		ID:                   domain.ID,
		Proposal:             responseForProposal,
		Name:                 domain.Name,
		EstimatedHarvestDate: domain.EstimatedHarvestDate,
		Status:               domain.Status,
		CancelReason:         domain.CancelReason,
		IsAvailable:          domain.IsAvailable,
		CreatedAt:            domain.CreatedAt,
		UpdatedAt:            domain.UpdatedAt,
	}, http.StatusOK, nil
}

func FromDomainArray(domain []batchs.Domain, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) ([]Batch, int, error) {
	var batches []Batch
	for _, value := range domain {
		batch, statusCode, err := FromDomain(value, proposalUC, commodityUC, userUC, regionUC)
		if err != nil {
			return []Batch{}, statusCode, err
		}

		batches = append(batches, batch)
	}

	return batches, http.StatusOK, nil
}

type BatchWithoutProposal struct {
	ID                   primitive.ObjectID `json:"_id"`
	Name                 string             `json:"name"`
	EstimatedHarvestDate primitive.DateTime `json:"estimatedHarvestDate"`
	Status               string             `json:"status"`
	CancelReason         string             `json:"cancelReason,omitempty"`
	IsAvailable          bool               `json:"isAvailable"`
	CreatedAt            primitive.DateTime `json:"createdAt"`
	UpdatedAt            primitive.DateTime `json:"updatedAt,omitempty"`
}

func FromDomainWithoutProposal(domain *batchs.Domain) BatchWithoutProposal {
	return BatchWithoutProposal{
		ID:                   domain.ID,
		Name:                 domain.Name,
		EstimatedHarvestDate: domain.EstimatedHarvestDate,
		Status:               domain.Status,
		CancelReason:         domain.CancelReason,
		IsAvailable:          domain.IsAvailable,
		CreatedAt:            domain.CreatedAt,
		UpdatedAt:            domain.UpdatedAt,
	}
}
