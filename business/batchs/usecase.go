package batchs

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

type BatchUseCase struct {
	batchRepository     Repository
	proposalRepository  proposals.Repository
	commodityRepository commodities.Repository
}

func NewUseCase(br Repository, pr proposals.Repository, cr commodities.Repository) UseCase {
	return &BatchUseCase{
		batchRepository:     br,
		proposalRepository:  pr,
		commodityRepository: cr,
	}
}

/*
Create
*/

func (bu *BatchUseCase) Create(proposalID primitive.ObjectID, farmerID primitive.ObjectID) (int, error) {
	proposal, err := bu.proposalRepository.GetByID(proposalID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mendapatkan proposal")
	}

	if proposal.Status != constant.ProposalStatusApproved {
		return http.StatusBadRequest, errors.New("proposal belum disetujui")
	} else if !proposal.IsAvailable {
		return http.StatusBadRequest, errors.New("proposal tidak tersedia")
	}

	commodity, err := bu.commodityRepository.GetByIDWithoutDeleted(proposal.CommodityID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mendapatkan periode tanam")
	}

	if farmerID != commodity.FarmerID {
		return http.StatusForbidden, errors.New("proposal tidak ditemukan")
	}

	if !commodity.IsPerennials {
		return http.StatusBadRequest, errors.New("komoditas ini tidak bisa dibuat batch")
	}

	proposal.IsAvailable = false
	proposal.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = bu.proposalRepository.Update(&proposal)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengubah proposal")
	}

	lastBatch, err := bu.batchRepository.CountByProposalCode(proposal.Code)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	domain := &Domain{
		ID:                   primitive.NewObjectID(),
		ProposalID:           proposalID,
		Name:                 fmt.Sprintf("%s - %d", proposal.Name, lastBatch+1),
		EstimatedHarvestDate: primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, commodity.PlantingPeriod)),
		Status:               constant.BatchStatusPlanting,
		IsAvailable:          true,
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
	commodity, err := bu.commodityRepository.GetByID(commodityID)
	if err == mongo.ErrNoDocuments {
		return []Domain{}, http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return []Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
	}

	batchs, err := bu.batchRepository.GetByCommodityCode(commodity.Code)
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

func (bu *BatchUseCase) CountByYear(year int) (int, int, error) {
	statistic, err := bu.batchRepository.CountByYear(year)
	if err != nil {
		return 0, http.StatusInternalServerError, errors.New("gagal mendapatkan statistik")
	}

	return statistic, http.StatusOK, nil
}

func (bc *BatchUseCase) GetForTransactionByCommodityID(commodityID primitive.ObjectID) ([]Domain, int, error) {
	commodity, err := bc.commodityRepository.GetByID(commodityID)
	if err == mongo.ErrNoDocuments {
		return []Domain{}, http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return []Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
	}

	batchs, err := bc.batchRepository.GetForTransactionByCommodityCode(commodity.Code)
	if err != nil {
		return []Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	return batchs, http.StatusOK, nil
}

func (bc *BatchUseCase) GetForTransactionByID(id primitive.ObjectID) (Domain, int, error) {
	batch, err := bc.batchRepository.GetForTransactionByID(id)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	return batch, http.StatusOK, nil
}

func (bc *BatchUseCase) GetForHarvestByFarmerID(farmerID primitive.ObjectID) ([]Domain, int, error) {
	batchs, err := bc.batchRepository.GetForHarvestByFarmerID(farmerID)
	if err != nil {
		return []Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	return batchs, http.StatusOK, nil
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
