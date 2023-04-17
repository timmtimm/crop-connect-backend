package response

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/transactions"
	"crop_connect/business/users"
	"crop_connect/controller/transactions/response"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Batch struct {
	ID                   primitive.ObjectID `json:"_id"`
	Transaction          response.Buyer     `json:"transaction"`
	Name                 string             `json:"name"`
	EstimatedHarvestDate primitive.DateTime `json:"estimatedHarvestDate"`
	Status               string             `json:"status"`
	CancelReason         string             `json:"cancelReason,omitempty"`
	CreatedAt            primitive.DateTime `json:"createdAt"`
	UpdatedAt            primitive.DateTime `json:"updatedAt,omitempty"`
}

func FromDomain(domain batchs.Domain, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) (Batch, int, error) {
	transaction, statusCode, err := transactionUC.GetByID(domain.TransactionID)
	if err != nil {
		return Batch{}, statusCode, errors.New("gagal mendapatkan transaksi")
	}

	transactionResponse, statusCode, err := response.FromDomainToBuyer(&transaction, proposalUC, commodityUC, userUC, regionUC)
	if err != nil {
		return Batch{}, statusCode, errors.New("gagal mendapatkan transaksi")
	}

	return Batch{
		ID:                   domain.ID,
		Transaction:          transactionResponse,
		Name:                 domain.Name,
		EstimatedHarvestDate: domain.EstimatedHarvestDate,
		Status:               domain.Status,
		CancelReason:         domain.CancelReason,
		CreatedAt:            domain.CreatedAt,
		UpdatedAt:            domain.UpdatedAt,
	}, http.StatusOK, nil
}

func FromDomainArray(domain []batchs.Domain, transactionUC transactions.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) ([]Batch, int, error) {
	var batches []Batch
	for _, value := range domain {
		batch, statusCode, err := FromDomain(value, transactionUC, proposalUC, commodityUC, userUC, regionUC)
		if err != nil {
			return []Batch{}, statusCode, err
		}

		batches = append(batches, batch)
	}

	return batches, http.StatusOK, nil
}

type BatchWithoutTransaction struct {
	ID                   primitive.ObjectID `json:"_id"`
	Name                 string             `json:"name"`
	EstimatedHarvestDate primitive.DateTime `json:"estimatedHarvestDate"`
	Status               string             `json:"status"`
	CancelReason         string             `json:"cancelReason,omitempty"`
	CreatedAt            primitive.DateTime `json:"createdAt"`
	UpdatedAt            primitive.DateTime `json:"updatedAt,omitempty"`
}

func FromDomainWithoutTransaction(domain *batchs.Domain) BatchWithoutTransaction {
	return BatchWithoutTransaction{
		ID:                   domain.ID,
		Name:                 domain.Name,
		EstimatedHarvestDate: domain.EstimatedHarvestDate,
		Status:               domain.Status,
		CancelReason:         domain.CancelReason,
		CreatedAt:            domain.CreatedAt,
		UpdatedAt:            domain.UpdatedAt,
	}
}
