package harvests

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/transactions"
	treatmentRecords "crop_connect/business/treatment_records"
	"crop_connect/constant"
	"crop_connect/dto"
	"crop_connect/helper"
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

func (hu *HarvestUseCase) CheckFarmerIDByProposalID(proposalID primitive.ObjectID, farmerID primitive.ObjectID) (proposals.Domain, commodities.Domain, int, error) {
	proposal, err := hu.proposalRepository.GetByIDWithoutDeleted(proposalID)
	if err == mongo.ErrNoDocuments {
		return proposals.Domain{}, commodities.Domain{}, http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return proposals.Domain{}, commodities.Domain{}, http.StatusInternalServerError, errors.New("proposal tidak ditemukan")
	}

	commodity, err := hu.commodityRepository.GetByIDWithoutDeleted(proposal.CommodityID)
	if err == mongo.ErrNoDocuments {
		return proposals.Domain{}, commodities.Domain{}, http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return proposals.Domain{}, commodities.Domain{}, http.StatusInternalServerError, errors.New("komoditas tidak ditemukan")
	}

	if commodity.FarmerID != farmerID {
		return proposals.Domain{}, commodities.Domain{}, http.StatusForbidden, errors.New("anda tidak memiliki akses")
	}

	return proposal, commodity, http.StatusOK, nil
}

/*
Create
*/

func (hu *HarvestUseCase) SubmitHarvest(domain *Domain, farmerID primitive.ObjectID, images []*multipart.FileHeader, notes []string) (Domain, int, error) {
	checkBatch, err := hu.batchRepository.GetByID(domain.BatchID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	_, _, statusCode, err := hu.CheckFarmerIDByProposalID(checkBatch.ProposalID, farmerID)
	if err != nil {
		return Domain{}, statusCode, err
	}
	newestTreatmentRecord, err := hu.treatmentRecordRepository.GetNewestByBatchIDAndStatus(domain.BatchID, constant.TreatmentRecordStatusApproved)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("batch belum memiliki riwayat perawatan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan riwayat perawatan terbaru")
	}
	if newestTreatmentRecord.Date > domain.Date {
		return Domain{}, http.StatusBadRequest, errors.New("tanggal panen tidak boleh lebih awal dari tanggal perawatan terakhir")
	} else if domain.Date > primitive.NewDateTimeFromTime(time.Now()) {
		return Domain{}, http.StatusBadRequest, errors.New("tanggal panen tidak boleh lebih dari tanggal hari ini")
	}
	checkHarvest, err := hu.harvestRepository.GetByBatchIDAndStatus(domain.BatchID, "")
	if err == mongo.ErrNoDocuments {
		var imageURLs []string

		if len(images) > 0 && len(notes) > 0 {

			imageURLs, err = hu.cloudinary.UploadManyWithGeneratedFilename(constant.CloudinaryFolderHarvests, images)
			if err != nil {
				return Domain{}, http.StatusInternalServerError, errors.New("gagal mengunggah gambar")
			}

			tempImageAndNotes := []dto.ImageAndNote{}
			for i := 0; i < len(imageURLs); i++ {
				tempImageAndNotes = append(tempImageAndNotes, dto.ImageAndNote{
					ImageURL: imageURLs[i],
					Note:     notes[i],
				})
			}

			domain.Harvest = tempImageAndNotes
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

func (hu *HarvestUseCase) GetByBatchIDAndStatus(batchID primitive.ObjectID, status string) (Domain, int, error) {
	harvest, err := hu.harvestRepository.GetByBatchIDAndStatus(batchID, status)
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

func (hu *HarvestUseCase) CountByYear(year int) (float64, int, error) {
	count, err := hu.harvestRepository.CountByYear(year)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return count, http.StatusOK, nil
}

func (hu *HarvestUseCase) GetByID(id primitive.ObjectID) (Domain, int, error) {
	harvest, err := hu.harvestRepository.GetByID(id)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("hasil panen tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan hasil panen")
	}

	return harvest, http.StatusOK, nil
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
		domain.RevisionNote = ""

		batch, err := hu.batchRepository.GetByID(harvest.BatchID)
		if err == mongo.ErrNoDocuments {
			return Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
		} else if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
		}

		batch.Status = constant.BatchStatusHarvest
		batch.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = hu.batchRepository.Update(&batch)
		if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal memperbarui batch")
		}

		proposal, err := hu.proposalRepository.GetByID(batch.ProposalID)
		if err == mongo.ErrNoDocuments {
			return Domain{}, http.StatusNotFound, errors.New("proposal tidak ditemukan")
		} else if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan proposal")
		}

		proposal.IsAvailable = true
		proposal.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = hu.proposalRepository.Update(&proposal)
		if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal memperbarui proposal")
		}
	}

	if domain.Status == constant.HarvestStatusRevision {
		if domain.RevisionNote == "" {
			return Domain{}, http.StatusBadRequest, errors.New("catatan revisi tidak boleh kosong")
		}

		harvest.RevisionNote = domain.RevisionNote
	}

	harvest.Status = domain.Status
	harvest.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return *domain, http.StatusOK, nil
}

func (hu *HarvestUseCase) UpdateHarvest(domain *Domain, farmerID primitive.ObjectID, updateImages []*helper.UpdateImage, notes []string) (Domain, int, error) {
	harvest, err := hu.harvestRepository.GetByID(domain.ID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("panen tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan panen")
	}

	if harvest.Status == constant.HarvestStatusApproved {
		return Domain{}, http.StatusConflict, errors.New("panen sudah diterima")
	}

	batch, err := hu.batchRepository.GetByID(harvest.BatchID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("batch tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan batch")
	}

	_, _, statusCode, err := hu.CheckFarmerIDByProposalID(batch.ProposalID, farmerID)
	if err != nil {
		return Domain{}, statusCode, err
	}

	if len(updateImages) > 0 && len(notes) > 0 {
		imageURLs := []string{}
		for _, imageAndNote := range harvest.Harvest {
			imageURLs = append(imageURLs, imageAndNote.ImageURL)
		}

		newImageURLs, err := hu.cloudinary.UpdateArrayImage(constant.CloudinaryFolderCommodities, imageURLs, updateImages)
		if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mengupdate gambar")
		}

		for i := 0; i < len(newImageURLs); i++ {
			if len(newImageURLs) == i {
				harvest.Harvest = append(harvest.Harvest, dto.ImageAndNote{
					ImageURL: newImageURLs[i],
					Note:     notes[i],
				})
			} else {
				harvest.Harvest[i] = dto.ImageAndNote{
					ImageURL: newImageURLs[i],
					Note:     notes[i],
				}
			}
		}
	} else {
		return Domain{}, http.StatusBadRequest, errors.New("gambar dan catatan tidak boleh kosong")
	}

	harvest.Status = constant.HarvestStatusPending
	harvest.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	harvest, err = hu.harvestRepository.Update(&harvest)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal memperbarui panen")
	}

	return harvest, http.StatusOK, nil
}

/*
Delete
*/
