package transactions

import (
	"errors"
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/proposals"
	"marketplace-backend/constant"
	"net/http"
	"time"

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

func (tu *TransactionUseCase) Create(domain *Domain) (int, error) {
	_, err := tu.transactionRepository.GetByBuyerIDProposalIDAndStatus(domain.BuyerID, domain.ProposalID, constant.TransactionStatusPending)
	if err == mongo.ErrNoDocuments {
		proposal, err := tu.proposalRepository.GetByID(domain.ProposalID)
		if err == mongo.ErrNoDocuments {
			return http.StatusNotFound, errors.New("proposal tidak ditemukan")
		} else if err != nil {
			return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
		}

		commodity, err := tu.commodityRepository.GetByID(proposal.CommodityID)
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
		domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = tu.transactionRepository.Create(domain)
		if err != nil {
			return http.StatusInternalServerError, errors.New("gagal membuat transaksi")
		}

		return http.StatusOK, nil
	} else {
		return http.StatusConflict, errors.New("transaksi sedang diproses")
	}

}

/*
Read
*/

func (tu *TransactionUseCase) GetByPaginationAndQuery(query Query) ([]Domain, int, int, error) {
	commodities, totalData, err := tu.transactionRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, err
	}

	return commodities, totalData, http.StatusOK, nil
}

func (tu *TransactionUseCase) GetTransactionsByCommodityName(query Query) ([]Domain, int, int, error) {
	commodities, totalData, err := tu.transactionRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
	}

	return commodities, totalData, http.StatusOK, nil
}

/*
Update
*/

func (tu *TransactionUseCase) MakeDecision(domain *Domain) (int, error) {
	transactions, err := tu.transactionRepository.GetByIDAndFarmerID(domain.ID, domain.BuyerID)
	if err != nil {
		return http.StatusNotFound, errors.New("transaksi tidak ditemukan")
	}

	if transactions.Status != constant.TransactionStatusPending {
		return http.StatusConflict, errors.New("transaksi sudah dibuat keputusan")
	}

	if domain.Status == constant.TransactionStatusAccepted {
		proposal, err := tu.proposalRepository.GetByID(transactions.ProposalID)
		if err != nil {
			return http.StatusNotFound, errors.New("proposal tidak ditemukan")
		}

		err = tu.transactionRepository.RejectPendingByProposalID(transactions.ProposalID)
		if err != nil {
			return http.StatusInternalServerError, errors.New("gagal mengupdate transaksi")
		}

		proposal.IsAvailable = false
		proposal.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = tu.proposalRepository.Update(&proposal)
		if err != nil {
			return http.StatusInternalServerError, errors.New("gagal mengupdate proposal")
		}
	}

	transactions.Status = domain.Status
	transactions.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = tu.transactionRepository.Update(&transactions)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengupdate transaksi")
	}

	return http.StatusOK, nil
}

/*
Delete
*/
