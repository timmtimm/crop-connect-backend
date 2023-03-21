package batchs

import (
	"errors"
	"fmt"
	"marketplace-backend/business/commodities"
	"marketplace-backend/business/proposals"
	"marketplace-backend/business/transactions"
	"marketplace-backend/constant"
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

func NewBatchUseCase(br Repository, tr transactions.Repository, pr proposals.Repository, cr commodities.Repository) UseCase {
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

/*
Update
*/

/*
Delete
*/
