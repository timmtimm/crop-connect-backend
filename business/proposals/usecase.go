package proposals

import (
	"errors"
	"marketplace-backend/business/commodities"
	"marketplace-backend/constant"
	"marketplace-backend/util"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProposalUseCase struct {
	proposalRepository  Repository
	commodityRepository commodities.Repository
}

func NewProposalUseCase(pr Repository, cr commodities.Repository) UseCase {
	return &ProposalUseCase{
		proposalRepository:  pr,
		commodityRepository: cr,
	}
}

/*
Create
*/

func (pu *ProposalUseCase) Create(domain *Domain, farmerID primitive.ObjectID) (int, error) {
	_, err := pu.commodityRepository.GetByIDAndFarmerID(domain.CommodityID, farmerID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	_, err = pu.proposalRepository.GetByCommodityIDAndName(domain.CommodityID, domain.Name)
	if err == mongo.ErrNoDocuments {
		domain.ID = primitive.NewObjectID()
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

		_, err = pu.proposalRepository.Update(&proposal)
		if err != nil {
			return http.StatusInternalServerError, errors.New("gagal memperbarui proposal")
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
