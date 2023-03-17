package transactions

import (
	"errors"
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/proposals"
	"marketplace-backend/constant"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionUseCase struct {
	transactionRepository Repository
	commodityRepository   commodities.Repository
	proposalRepository    proposals.Repository
}

func NewTransactionUseCase(tr Repository, cr commodities.Repository, pr proposals.Repository) UseCase {
	return &TransactionUseCase{
		transactionRepository: tr,
		commodityRepository:   cr,
		proposalRepository:    pr,
	}
}

/*
Create
*/

func (tc *TransactionUseCase) Create(domain *Domain) (int, error) {

	proposal, err := tc.proposalRepository.GetByID(domain.ProposalID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	commodity, err := tc.commodityRepository.GetByID(proposal.CommodityID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data komoditas")
	}

	if !proposal.IsAvailable {
		return http.StatusConflict, errors.New("proposal tidak tersedia")
	}

	domain.ID = primitive.NewObjectID()
	domain.Status = constant.TransactionStatusPending
	domain.TotalPrice = float64(commodity.PricePerKg) * proposal.EstimatedTotalHarvest

	_, err = tc.transactionRepository.Create(domain)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

/*
Read
*/

/*
Update
*/

/*
Delete
*/
