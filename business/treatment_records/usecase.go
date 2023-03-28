package treatment_records

import (
	"errors"
	"fmt"
	"marketplace-backend/business/batchs"
	"marketplace-backend/constant"
	"marketplace-backend/helper/cloudinary"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TreatmentRecordUseCase struct {
	treatmentRecordRepository Repository
	batchRepository           batchs.Repository
	cloudinary                cloudinary.Function
}

func NewTreatmentRecordUseCase(trr Repository, br batchs.Repository, cldry cloudinary.Function) UseCase {
	return &TreatmentRecordUseCase{
		treatmentRecordRepository: trr,
		batchRepository:           br,
		cloudinary:                cldry,
	}
}

/*
Create
*/

func (tru *TreatmentRecordUseCase) RequestToFarmer(domain *Domain) (Domain, int, error) {
	batch, err := tru.batchRepository.GetByID(domain.BatchID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	if batch.Status != constant.BatchStatusPlanting {
		return Domain{}, http.StatusBadRequest, errors.New("batch tidak sedang dalam tahap tanam")
	} else if batch.CreatedAt >= domain.Date {
		return Domain{}, http.StatusBadRequest, errors.New("tanggal perawatan harus lebih besar dari tanggal tanam")
	}

	count, err := tru.treatmentRecordRepository.CountByBatchID(domain.BatchID)
	fmt.Println(count)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan jumlah riwayat perawatan")
	}

	newestTreatmentRecord, err := tru.treatmentRecordRepository.GetNewestByBatchID(domain.BatchID)
	if err != mongo.ErrNoDocuments {
		if newestTreatmentRecord.Status != constant.TreatmentRecordStatusAccepted {
			return Domain{}, http.StatusBadRequest, errors.New("riwayat perawatan terakhir belum selesai")
		} else if primitive.NewDateTimeFromTime(time.Now()) >= domain.Date {
			return Domain{}, http.StatusBadRequest, errors.New("tanggal perawatan harus lebih besar dari tanggal hari ini")
		} else if newestTreatmentRecord.Date >= domain.Date {
			return Domain{}, http.StatusBadRequest, errors.New("tanggal perawatan harus lebih besar dari tanggal perawatan terakhir")
		}
	} else if err != mongo.ErrNoDocuments && err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan terakhir")
	}

	domain.ID = primitive.NewObjectID()
	domain.Number = count + 1
	domain.Status = constant.TreatmentRecordStatusWaitingResponse
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	treatmentRecord, err := tru.treatmentRecordRepository.Create(domain)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal membuat riwayat perawatan")
	}

	return treatmentRecord, http.StatusCreated, nil
}

/*
Read
*/

func (tru *TreatmentRecordUseCase) GetByPaginationAndQuery(query Query) ([]Domain, int, int, error) {
	treatmentRecords, totalData, err := tru.treatmentRecordRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan")
	}

	return treatmentRecords, totalData, http.StatusOK, nil
}

/*
Update
*/

/*
Delete
*/
