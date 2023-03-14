package proposals

import (
	"errors"
	"marketplace-backend/business/commodities"
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

func (pr *ProposalUseCase) Create(domain *Domain, farmerID primitive.ObjectID) (int, error) {
	_, err := pr.commodityRepository.GetByIDAndFarmerID(domain.CommodityID, farmerID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	_, err = pr.proposalRepository.GetByCommodityIDAndName(domain.CommodityID, domain.Name)
	if err == nil {
		return http.StatusConflict, errors.New("nama proposal sudah digunakan")
	}

	domain.ID = primitive.NewObjectID()
	domain.IsAccepted = false
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = pr.proposalRepository.Create(domain)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

/*
Read
*/

func (pr *ProposalUseCase) GetByCommodityID(commodityID primitive.ObjectID) ([]Domain, int, error) {
	proposals, err := pr.proposalRepository.GetByCommodityIDAndAvailability(commodityID, true)
	if err == mongo.ErrNoDocuments {
		return proposals, http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return proposals, http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	return proposals, http.StatusOK, nil
}

/*
Update
*/

func (pr *ProposalUseCase) Update(domain *Domain, farmerID primitive.ObjectID) (int, error) {
	proposal, err := pr.proposalRepository.GetByID(domain.ID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	_, err = pr.commodityRepository.GetByIDAndFarmerID(proposal.CommodityID, farmerID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	}

	if proposal.Name != domain.Name {
		_, err = pr.proposalRepository.GetByCommodityIDAndName(domain.CommodityID, domain.Name)
		if err == nil {
			return http.StatusConflict, errors.New("nama proposal sudah digunakan")
		}
	}

	if proposal.IsAccepted {
		err = pr.proposalRepository.Delete(proposal.ID)
		if err == mongo.ErrNilDocument {
			return http.StatusNotFound, errors.New("proposal tidak ditemukan")
		} else if err != nil {
			return http.StatusInternalServerError, errors.New("gagal menghapus proposal")
		}

		domain.ID = primitive.NewObjectID()
		domain.CommodityID = proposal.CommodityID
		domain.IsAccepted = false
		domain.CreatedAt = proposal.CreatedAt
		domain.UpdatedAt = proposal.UpdatedAt

		_, err = pr.proposalRepository.Create(domain)
		if err != nil {
			return http.StatusInternalServerError, errors.New("gagal membuat proposal")
		}

		return http.StatusOK, nil
	} else {
		proposal.Name = domain.Name
		proposal.Description = domain.Description
		proposal.EstimatedTotalHarvest = domain.EstimatedTotalHarvest
		proposal.PlantingArea = domain.PlantingArea
		proposal.Address = domain.Address
		proposal.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = pr.proposalRepository.Update(&proposal)
		if err == mongo.ErrNilDocument {
			return http.StatusNotFound, errors.New("proposal tidak ditemukan")
		} else if err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, nil
	}
}

func (pr *ProposalUseCase) UpdateCommodityID(oldCommodityID primitive.ObjectID, NewCommodityID primitive.ObjectID) (int, error) {
	proposals, err := pr.proposalRepository.GetByCommodityID(oldCommodityID)
	if err == mongo.ErrNoDocuments {
		return http.StatusOK, nil
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	for _, proposal := range proposals {
		proposal.CommodityID = NewCommodityID
		proposal.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = pr.proposalRepository.Update(&proposal)
		if err == mongo.ErrNilDocument {
			return http.StatusNotFound, errors.New("proposal tidak ditemukan")
		} else if err != nil {
			return http.StatusInternalServerError, errors.New("gagal memperbarui proposal")
		}
	}

	return http.StatusOK, nil
}

/*
Delete
*/

func (pr *ProposalUseCase) Delete(id primitive.ObjectID, farmerID primitive.ObjectID) (int, error) {
	proposal, err := pr.proposalRepository.GetByID(id)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	_, err = pr.commodityRepository.GetByIDAndFarmerID(proposal.CommodityID, farmerID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data komoditas")
	}

	err = pr.proposalRepository.Delete(id)
	if err == mongo.ErrNilDocument {
		return http.StatusNotFound, errors.New("proposal tidak ditemukan")
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (pr *ProposalUseCase) DeleteByCommodityID(commodityID primitive.ObjectID) (int, error) {
	proposals, err := pr.proposalRepository.GetByCommodityID(commodityID)
	if err == mongo.ErrNoDocuments {
		return http.StatusOK, nil
	} else if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	for _, proposal := range proposals {
		err = pr.proposalRepository.Delete(proposal.ID)
		if err == mongo.ErrNilDocument {
			return http.StatusNotFound, errors.New("proposal tidak ditemukan")
		} else if err != nil {
			return http.StatusInternalServerError, errors.New("gagal menghapus proposal")
		}
	}

	return http.StatusOK, nil
}
