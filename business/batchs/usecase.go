package batchs

import (
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/transactions"
	"crop_connect/constant"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BatchUseCase struct {
	batchRepository       Repository
	transactionRepository transactions.Repository
	proposalRepository    proposals.Repository
	commodityRepository   commodities.Repository
}

func NewUseCase(br Repository, tr transactions.Repository, pr proposals.Repository, cr commodities.Repository) UseCase {
	return &BatchUseCase{
		batchRepository:       br,
		transactionRepository: tr,
		proposalRepository:    pr,
		commodityRepository:   cr,
	}
}

/*
Create
*/

func (bu *BatchUseCase) Create(transactionID primitive.ObjectID) (int, error) {
	transaction, err := bu.transactionRepository.GetByID(transactionID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("transaksi tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mendapatkan transaksi")
	}

	proposal, err := bu.proposalRepository.GetByID(transaction.ProposalID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mendapatkan proposal")
	}

	lastBatch, err := bu.batchRepository.CountByProposalName(proposal.Name)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal menghitung jumlah batch")
	}

	commodity, err := bu.commodityRepository.GetByID(proposal.CommodityID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mendapatkan periode tanam")
	}

	domain := &Domain{
		ID:                   primitive.NewObjectID(),
		TransactionID:        transactionID,
		Name:                 fmt.Sprintf("%s - %d", proposal.Name, lastBatch+1),
		EstimatedHarvestDate: primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, commodity.PlantingPeriod)),
		Status:               constant.BatchStatusPlanting,
		CreatedAt:            primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err = bu.batchRepository.Create(domain)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal membuat batch")
	}

	return http.StatusCreated, nil
}

/*
Read
*/

func (bu *BatchUseCase) GetByID(id primitive.ObjectID) (Domain, int, error) {
	batch, err := bu.batchRepository.GetByID(id)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	return batch, http.StatusOK, nil
}

func (bu *BatchUseCase) GetByCommodityID(commodityID primitive.ObjectID) ([]Domain, int, error) {
	batchs, err := bu.batchRepository.GetByCommodityID(commodityID)
	if err != nil {
		return []Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	return batchs, http.StatusOK, nil
}

func (bu *BatchUseCase) GetByPaginationAndQuery(query Query) ([]Domain, int, int, error) {
	batches, totalData, err := bu.batchRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, 0, errors.New("gagal mendapatkan batch")
	}

	return batches, totalData, http.StatusOK, nil
}

/*
Update
*/

// func (bu *BatchUseCase) Cancel(domain *Domain, farmerID primitive.ObjectID) (int, error) {
// 	batch, err := bu.batchRepository.GetByID(domain.ID)
// 	if err == mongo.ErrNoDocuments {
// 		return http.StatusNotFound, errors.New("batch tidak ditemukan")
// 	} else if err != nil {
// 		return http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
// 	}

// 	transaction, err := bu.transactionRepository.GetByID(batch.TransactionID)
// 	if err == mongo.ErrNoDocuments {
// 		return http.StatusNotFound, errors.New("transaksi tidak ditemukan")
// 	} else if err != nil {
// 		return http.StatusInternalServerError, errors.New("gagal mendapatkan transaksi")
// 	}

// 	transaction.Status = constant.TransactionStatusCancel
// 	_, err = bu.transactionRepository.Update(&transaction)
// 	if err != nil {
// 		return http.StatusInternalServerError, errors.New("gagal membatalkan transaksi")
// 	}

// 	proposal, err := bu.proposalRepository.GetByID(transaction.ProposalID)
// 	if err == mongo.ErrNoDocuments {
// 		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
// 	} else if err != nil {
// 		return http.StatusInternalServerError, errors.New("gagal mendapatkan proposal")
// 	}

// 	commodity, err := bu.commodityRepository.GetByID(proposal.CommodityID)
// 	if err == mongo.ErrNoDocuments {
// 		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
// 	} else if err != nil {
// 		return http.StatusInternalServerError, errors.New("gagal mendapatkan periode tanam")
// 	}

// 	if commodity.FarmerID != farmerID {
// 		return http.StatusForbidden, errors.New("tidak memiliki akses")
// 	}

// 	if batch.Status == constant.BatchStatusCancel {
// 		return http.StatusBadRequest, errors.New("batch sudah dibatalkan")
// 	}

// 	batch.Status = constant.BatchStatusCancel
// 	batch.CancelReason = domain.CancelReason
// 	batch.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

// 	_, err = bu.batchRepository.Update(&batch)
// 	if err != nil {
// 		return http.StatusInternalServerError, errors.New("gagal membatalkan batch")
// 	}

// 	return http.StatusOK, nil
// }

/*
Delete
*/
