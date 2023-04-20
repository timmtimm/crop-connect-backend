package harvests

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/transactions"
	treatmentRecords "crop_connect/business/treatment_records"
	"crop_connect/constant"
	"crop_connect/dto"
	"crop_connect/helper/cloudinary"
	"crop_connect/util"
	"errors"
	"mime/multipart"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HarvestUseCase struct {
	harvestRepository         Repository
	treatmentRecordRepository treatmentRecords.Repository
	batchRepository           batchs.Repository
	transactionRepository     transactions.Repository
	proposalRepository        proposals.Repository
	commodityRepository       commodities.Repository
	cloudinary                cloudinary.Function
}

func NewUseCase(hr Repository, br batchs.Repository, trr treatmentRecords.Repository, tr transactions.Repository, pr proposals.Repository, cr commodities.Repository, cldry cloudinary.Function) UseCase {
	return &HarvestUseCase{
		harvestRepository:         hr,
		treatmentRecordRepository: trr,
		batchRepository:           br,
		transactionRepository:     tr,
		proposalRepository:        pr,
		commodityRepository:       cr,
		cloudinary:                cldry,
	}
}

/*
Create
*/

func (hu *HarvestUseCase) SubmitHarvest(domain *Domain, farmerID primitive.ObjectID, images []*multipart.FileHeader, notes []string) (Domain, int, error) {
	newestTreatmentRecord, err := hu.treatmentRecordRepository.GetNewestByBatchID(domain.BatchID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("riwayat perawatan tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan")
	}

	if newestTreatmentRecord.Date > domain.Date {
		return Domain{}, http.StatusBadRequest, errors.New("tanggal panen tidak boleh lebih awal dari tanggal perawatan terakhir")
	} else if domain.Date > primitive.NewDateTimeFromTime(time.Now()) {
		return Domain{}, http.StatusBadRequest, errors.New("tanggal panen tidak boleh lebih dari tanggal hari ini")
	}

	batch, err := hu.batchRepository.GetByID(domain.BatchID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	checkHarvest, err := hu.harvestRepository.GetByBatchID(domain.BatchID)
	if err == mongo.ErrNoDocuments {
		transaction, err := hu.transactionRepository.GetByID(batch.TransactionID)
		if err == mongo.ErrNoDocuments {
			return Domain{}, http.StatusNotFound, errors.New("transaksi tidak ditemukan")
		} else if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan transaksi")
		}

		proposal, err := hu.proposalRepository.GetByID(transaction.ProposalID)
		if err == mongo.ErrNoDocuments {
			return Domain{}, http.StatusNotFound, errors.New("proposal tidak ditemukan")
		} else if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan proposal")
		}

		commodity, err := hu.commodityRepository.GetByID(proposal.CommodityID)
		if err == mongo.ErrNoDocuments {
			return Domain{}, http.StatusNotFound, errors.New("komoditas tidak ditemukan")
		} else if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
		}

		if commodity.FarmerID != farmerID {
			return Domain{}, http.StatusForbidden, errors.New("anda tidak memiliki akses")
		}

		var imageURLs []string
		notes = util.RemoveNilStringInArray(notes)

		if len(images) != len(notes) {
			return Domain{}, http.StatusBadRequest, errors.New("jumlah gambar dan catatan tidak sama")
		}

		if len(images) > 0 {
			imageURLs, err = hu.cloudinary.UploadManyWithGeneratedFilename(constant.CloudinaryFolderHarvests, images)
			if err != nil {
				return Domain{}, http.StatusInternalServerError, errors.New("gagal mengunggah gambar")
			}

			for i := 0; i < len(imageURLs); i++ {
				domain.Harvest = append(domain.Harvest, dto.ImageAndNote{
					ImageURL: imageURLs[i],
					Note:     notes[i],
				})
			}
		} else {
			return Domain{}, http.StatusBadRequest, errors.New("gambar dan catatan tidak boleh kosong")
		}

		domain.ID = primitive.NewObjectID()
		domain.Status = constant.HarvestStatusPending
		domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = hu.harvestRepository.Create(domain)
		if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mengajukan hasi panen")
		}

		return *domain, http.StatusCreated, nil
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan hasil panen")
	}

	if checkHarvest.Status == constant.HarvestStatusPending {
		return Domain{}, http.StatusBadRequest, errors.New("hasil panen sedang dalam proses verifikasi")
	} else if checkHarvest.Status == constant.HarvestStatusApproved {
		return Domain{}, http.StatusBadRequest, errors.New("hasil panen sudah diterima")
	} else {
		return Domain{}, http.StatusBadRequest, errors.New("hasil panen sedang dalam proses revisi")
	}
}

/*
Read
*/

func (hu *HarvestUseCase) GetByBatchID(batchID primitive.ObjectID) (Domain, int, error) {
	harvest, err := hu.harvestRepository.GetByBatchID(batchID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("hasil panen tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan hasil panen")
	}

	return harvest, http.StatusOK, nil
}

func (hu *HarvestUseCase) GetByPaginationAndQuery(query Query) ([]Domain, int, int, error) {
	harvests, totalData, err := hu.harvestRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, errors.New("gagal mendapatkan hasil panen")
	}

	return harvests, totalData, http.StatusOK, nil
}

/*
Update
*/

func (hu *HarvestUseCase) Validate(domain *Domain, validatorID primitive.ObjectID) (Domain, int, error) {
	isStatusAvailable := util.CheckStringOnArray([]string{constant.HarvestStatusRevision, constant.HarvestStatusApproved}, domain.Status)
	if !isStatusAvailable {
		return Domain{}, http.StatusBadRequest, errors.New("status harvest hanya tersedia approved dan revision")
	}

	harvest, err := hu.harvestRepository.GetByID(domain.ID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("hasil panen tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan hasil panen")
	}

	if harvest.Status != constant.HarvestStatusPending {
		return Domain{}, http.StatusBadRequest, errors.New("hasil panen tidak sedang dalam proses verifikasi")
	}

	if domain.Status == constant.HarvestStatusApproved {
		domain.AccepterID = validatorID
	}

	if domain.Status == constant.HarvestStatusRevision {
		if domain.RevisionNote == "" {
			return Domain{}, http.StatusBadRequest, errors.New("catatan revisi tidak boleh kosong")
		}

		harvest.RevisionNote = domain.RevisionNote
	}

	harvest.Status = domain.Status
	harvest.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = hu.harvestRepository.Update(&harvest)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal memvalidasi hasil panen")
	}

	return *domain, http.StatusOK, nil
}

/*
Delete
*/
