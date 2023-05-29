package proposals

import (
	"crop_connect/business/commodities"
	"crop_connect/business/regions"
	"crop_connect/constant"
	"crop_connect/dto"
	"crop_connect/util"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProposalUseCase struct {
	proposalRepository  Repository
	commodityRepository commodities.Repository
	regionRepository    regions.Repository
}

func NewUseCase(pr Repository, cr commodities.Repository, rr regions.Repository) UseCase {
	return &ProposalUseCase{
		proposalRepository:  pr,
		commodityRepository: cr,
		regionRepository:    rr,
	}
}

/*
Create
*/

func (pu *ProposalUseCase) Create(domain *Domain, farmerID primitive.ObjectID) (int, error) {
	_, err := pu.regionRepository.GetByID(domain.RegionID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("daerah tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data daerah")
	}

	_, err = pu.commodityRepository.GetByIDAndFarmerID(domain.CommodityID, farmerID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data komoditas")
	}

	_, err = pu.proposalRepository.GetByCommodityIDAndName(domain.CommodityID, domain.Name)
	if err == mongo.ErrNoDocuments {
		domain.ID = primitive.NewObjectID()
		domain.Code = primitive.NewObjectID()
		domain.Status = constant.ProposalStatusPending
		domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = pu.proposalRepository.Create(domain)
		if err != nil {
			return http.StatusInternalServerError, errors.New("gagal membuat proposal")
		}

		return http.StatusCreated, nil
	} else {
		return http.StatusConflict, errors.New("nama proposal sudah digunakan")
	}
}

/*
Read
*/

func (pu *ProposalUseCase) GetByID(id primitive.ObjectID) (Domain, int, error) {
	proposal, err := pu.proposalRepository.GetByID(id)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	return proposal, http.StatusOK, nil
}

func (pu *ProposalUseCase) GetByCommodityID(commodityID primitive.ObjectID) ([]Domain, int, error) {
	proposals, err := pu.proposalRepository.GetByCommodityIDAndAvailability(commodityID, constant.ProposalStatusApproved)
	if err == mongo.ErrNoDocuments {
		return []Domain{}, http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return []Domain{}, http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	return proposals, http.StatusOK, nil
}

func (pu *ProposalUseCase) GetByIDWithoutDeleted(id primitive.ObjectID) (Domain, int, error) {
	proposals, err := pu.proposalRepository.GetByIDWithoutDeleted(id)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	return proposals, http.StatusOK, nil
}

func (pu *ProposalUseCase) GetByIDAccepted(id primitive.ObjectID) (Domain, int, error) {
	proposal, err := pu.proposalRepository.GetByIDAccepted(id)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mendapatkan proposal")
	}

	return proposal, http.StatusOK, nil
}

func (pu *ProposalUseCase) StatisticByYear(year int) ([]dto.StatisticByYear, int, error) {
	totalProposal, err := pu.proposalRepository.StatisticByYear(year)
	if err != nil {
		return []dto.StatisticByYear{}, http.StatusInternalServerError, err
	}

	if len(totalProposal) < 12 {
		totalProposal = util.FillNotAvailableMonth(totalProposal)
	}

	return totalProposal, http.StatusOK, nil
}

func (pu *ProposalUseCase) CountTotalProposalByFarmer(farmerID primitive.ObjectID) (int, int, error) {
	totalProposal, err := pu.proposalRepository.CountTotalProposalByFarmer(farmerID)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return totalProposal, http.StatusOK, nil
}

func (pu *ProposalUseCase) GetByPaginationAndQuery(query Query) ([]Domain, int, int, error) {
	proposals, total, err := pu.proposalRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	return proposals, total, http.StatusOK, nil
}

/*
Update
*/

func (pu *ProposalUseCase) Update(domain *Domain, farmerID primitive.ObjectID) (int, error) {
	proposal, err := pu.proposalRepository.GetByID(domain.ID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	_, err = pu.commodityRepository.GetByIDAndFarmerID(proposal.CommodityID, farmerID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	}

	if proposal.RegionID != domain.RegionID {
		_, err = pu.regionRepository.GetByID(domain.RegionID)
		if err == mongo.ErrNoDocuments {
			return http.StatusNotFound, errors.New("daerah tidak ditemukan")
		} else if err != nil {
			return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
		}
	}

	if proposal.Name != domain.Name {
		_, err = pu.proposalRepository.GetByCommodityIDAndName(domain.CommodityID, domain.Name)
		if err != mongo.ErrNoDocuments {
			return http.StatusConflict, errors.New("nama proposal sudah digunakan")
		}
	}

	if proposal.Status == constant.ProposalStatusApproved {
		err = pu.proposalRepository.Delete(proposal.ID)
		if err == mongo.ErrNoDocuments {
			return http.StatusNotFound, errors.New("proposal tidak ditemukan")
		} else if err != nil {
			return http.StatusInternalServerError, errors.New("gagal menghapus proposal")
		}

		domain.ID = primitive.NewObjectID()
		domain.CommodityID = proposal.CommodityID
		domain.Status = constant.ProposalStatusPending
		domain.CreatedAt = proposal.CreatedAt
		domain.UpdatedAt = proposal.UpdatedAt

		_, err = pu.proposalRepository.Create(domain)
		if err != nil {
			return http.StatusInternalServerError, errors.New("gagal membuat proposal")
		}
	} else if proposal.Status == constant.ProposalStatusPending || proposal.Status == constant.ProposalStatusRejected {
		proposal.Name = domain.Name
		proposal.Description = domain.Description
		proposal.Status = constant.ProposalStatusPending
		proposal.EstimatedTotalHarvest = domain.EstimatedTotalHarvest
		proposal.PlantingArea = domain.PlantingArea
		proposal.Address = domain.Address
		proposal.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = pu.proposalRepository.Update(&proposal)
		if err != nil {
			return http.StatusInternalServerError, errors.New("gagal memperbarui proposal")
		}
	} else {
		return http.StatusBadRequest, errors.New("status proposal tidak valid")
	}

	return http.StatusOK, nil
}

func (pu *ProposalUseCase) UpdateCommodityID(oldCommodityID primitive.ObjectID, NewCommodityID primitive.ObjectID) (int, error) {
	proposals, err := pu.proposalRepository.GetByCommodityID(oldCommodityID)
	if err == mongo.ErrNoDocuments {
		return http.StatusOK, nil
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	for _, proposal := range proposals {
		proposal.CommodityID = NewCommodityID
		proposal.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		if proposal.Status == constant.ProposalStatusApproved {
			err = pu.proposalRepository.Delete(proposal.ID)
			if err == mongo.ErrNoDocuments {
				return http.StatusNotFound, errors.New("proposal tidak ditemukan")
			} else if err != nil {
				return http.StatusInternalServerError, errors.New("gagal menghapus proposal")
			}

			proposal.ID = primitive.NewObjectID()
			_, err = pu.proposalRepository.Create(&proposal)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		} else if proposal.Status == constant.ProposalStatusPending || proposal.Status == constant.ProposalStatusRejected {
			_, err = pu.proposalRepository.Update(&proposal)
			if err != nil {
				return http.StatusInternalServerError, errors.New("gagal memperbarui proposal")
			}
		} else {
			return http.StatusBadRequest, errors.New("status proposal tidak valid")
		}
	}

	return http.StatusOK, nil
}

func (pu *ProposalUseCase) ValidateProposal(domain *Domain, validatorID primitive.ObjectID) (int, error) {
	proposal, err := pu.proposalRepository.GetByID(domain.ID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	if proposal.Status != constant.ProposalStatusPending {
		return http.StatusBadRequest, errors.New("proposal sudah divalidasi")
	}

	isStatusAvailable := util.CheckStringOnArray([]string{constant.ProposalStatusRejected, constant.ProposalStatusApproved}, domain.Status)
	if !isStatusAvailable {
		return http.StatusBadRequest, errors.New("status proposal hanya tersedia approved dan rejected")
	}

	proposal.ValidatorID = validatorID
	proposal.Status = domain.Status
	proposal.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	if domain.Status == constant.ProposalStatusRejected {
		proposal.RejectReason = domain.RejectReason
	} else {
		proposal.IsAvailable = true
		_, err = pu.proposalRepository.UnsetRejectReason(proposal.ID)
		if err != nil {
			return http.StatusInternalServerError, errors.New("gagal memperbarui proposal")
		}
	}

	_, err = pu.proposalRepository.Update(&proposal)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal memperbarui proposal")
	}

	return http.StatusOK, nil
}

/*
Delete
*/

func (pu *ProposalUseCase) Delete(id primitive.ObjectID, farmerID primitive.ObjectID) (int, error) {
	proposal, err := pu.proposalRepository.GetByID(id)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	_, err = pu.commodityRepository.GetByIDAndFarmerID(proposal.CommodityID, farmerID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data komoditas")
	}

	err = pu.proposalRepository.Delete(id)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal menghapus proposal")
	}

	return http.StatusOK, nil
}

func (pu *ProposalUseCase) DeleteByCommodityID(commodityID primitive.ObjectID) (int, error) {
	proposals, err := pu.proposalRepository.GetByCommodityID(commodityID)
	if err == mongo.ErrNoDocuments {
		return http.StatusOK, nil
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	for _, proposal := range proposals {
		err = pu.proposalRepository.Delete(proposal.ID)
		if err == mongo.ErrNoDocuments {
			return http.StatusNotFound, errors.New("proposal tidak ditemukan")
		} else if err != nil {
			return http.StatusInternalServerError, errors.New("gagal menghapus proposal")
		}
	}

	return http.StatusOK, nil
}
