package transactions

import (
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/constant"
	"errors"
	"fmt"
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

func NewUseCase(tr Repository, cr commodities.Repository, pr proposals.Repository) UseCase {
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

		return http.StatusCreated, nil
	} else {
		return http.StatusConflict, errors.New("transaksi sedang diproses")
	}

}

/*
Read
*/

func (tu *TransactionUseCase) GetByID(id primitive.ObjectID) (Domain, int, error) {
	transaction, err := tu.transactionRepository.GetByID(id)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("transaksi tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan transaksi")
	}

	return transaction, http.StatusOK, nil
}

func (tu *TransactionUseCase) GetByPaginationAndQuery(query Query) ([]Domain, int, int, error) {
	commodities, totalData, err := tu.transactionRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, errors.New("something")
	}

	return commodities, totalData, http.StatusOK, nil
}

func (tu *TransactionUseCase) GetByIDAndBuyerIDOrFarmerID(id primitive.ObjectID, buyerID primitive.ObjectID, farmerID primitive.ObjectID) (Domain, int, error) {
	transaction, err := tu.transactionRepository.GetByIDAndBuyerIDOrFarmerID(id, buyerID, farmerID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("transaksi tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan transaksi")
	}

	return transaction, http.StatusOK, nil
}

func (tu *TransactionUseCase) GetTransactionsByCommodityName(query Query) ([]Domain, int, int, error) {
	commodities, totalData, err := tu.transactionRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
	}

	return commodities, totalData, http.StatusOK, nil
}

func (tu *TransactionUseCase) StatisticByYear(farmerID primitive.ObjectID, year int) ([]Statistic, int, error) {
	statistics, err := tu.transactionRepository.StatisticByYear(farmerID, year)
	if err != nil {
		return []Statistic{}, http.StatusInternalServerError, err
	}

	return statistics, http.StatusOK, nil
}

func (tu *TransactionUseCase) StatisticTopProvince(year int, limit int) ([]TotalTransactionByProvince, int, error) {
	statistics, err := tu.transactionRepository.StatisticTopProvince(year, limit)
	if err != nil {
		return []TotalTransactionByProvince{}, http.StatusInternalServerError, err
	}

	return statistics, http.StatusOK, nil
}

func (tu *TransactionUseCase) StatisticTopCommodity(farmerID primitive.ObjectID, year int, limit int) ([]StatisticTopCommodity, int, error) {
	statistics, err := tu.transactionRepository.StatisticTopCommodity(farmerID, year, limit)
	if err != nil {
		return []StatisticTopCommodity{}, http.StatusInternalServerError, err
	}

	domainStatisticCommodity := []StatisticTopCommodity{}
	for _, statistic := range statistics {
		commodity, err := tu.commodityRepository.GetByCode(statistic.CommodityCode)
		if err != nil {
			return []StatisticTopCommodity{}, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
		}

		domainStatisticCommodity = append(domainStatisticCommodity, StatisticTopCommodity{
			Commodity: commodity,
			Total:     statistic.Total,
		})
	}

	return domainStatisticCommodity, http.StatusOK, nil
}

func (tu *TransactionUseCase) CountByCommodityID(commodityID primitive.ObjectID) (int, float64, int, error) {
	commodity, err := tu.commodityRepository.GetByID(commodityID)
	if err == mongo.ErrNoDocuments {
		return 0, 0, http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return 0, 0, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
	}

	fmt.Println(commodity)

	totalTransaction, totalWeight, err := tu.transactionRepository.CountByCommodityCode(commodity.Code)
	if err != nil {
		return 0, 0, http.StatusInternalServerError, err
	}

	return totalTransaction, totalWeight, http.StatusOK, nil
}

/*
Update
*/

func (tu *TransactionUseCase) MakeDecision(domain *Domain, farmerID primitive.ObjectID) (int, error) {
	transaction, err := tu.transactionRepository.GetByID(domain.ID)
	if err != nil {
		return http.StatusNotFound, errors.New("transaksi tidak ditemukan")
	}

	proposal, err := tu.proposalRepository.GetByIDWithoutDeleted(transaction.ProposalID)
	if err != nil {
		return http.StatusInternalServerError, errors.New("proposal tidak ditemukan")
	}

	commodity, err := tu.commodityRepository.GetByIDWithoutDeleted(proposal.CommodityID)
	if err != nil {
		return http.StatusInternalServerError, errors.New("komoditas tidak ditemukan")
	}

	if commodity.FarmerID != farmerID {
		return http.StatusForbidden, errors.New("anda tidak memiliki akses")
	}

	if transaction.Status != constant.TransactionStatusPending {
		return http.StatusConflict, errors.New("transaksi sudah dibuat keputusan")
	}

	if domain.Status == constant.TransactionStatusAccepted {
		proposal, err := tu.proposalRepository.GetByID(transaction.ProposalID)
		if err != nil {
			return http.StatusNotFound, errors.New("proposal tidak ditemukan")
		}

		err = tu.transactionRepository.RejectPendingByProposalID(transaction.ProposalID)
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

	transaction.Status = domain.Status
	transaction.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = tu.transactionRepository.Update(&transaction)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengupdate transaksi")
	}

	return http.StatusOK, nil
}

func (tu *TransactionUseCase) CancelOnPending(id primitive.ObjectID, buyerID primitive.ObjectID) (int, error) {
	transaction, err := tu.transactionRepository.GetByIDAndBuyerIDOrFarmerID(id, buyerID, primitive.NilObjectID)
	if err != nil {
		return http.StatusNotFound, errors.New("transaksi tidak ditemukan")
	}

	if transaction.Status != constant.TransactionStatusPending {
		return http.StatusConflict, errors.New("transaksi sudah dibuat keputusan")
	}

	transaction.Status = constant.TransactionStatusCancel
	transaction.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = tu.transactionRepository.Update(&transaction)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengupdate transaksi")
	}

	return http.StatusOK, nil
}

/*
Delete
*/
