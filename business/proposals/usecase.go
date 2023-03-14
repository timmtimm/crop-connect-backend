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

	commodity, err := pr.commodityRepository.GetByIDAndFarmerID(proposal.CommodityID, farmerID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	}

	if commodity.FarmerID != farmerID {
		return http.StatusForbidden, errors.New("anda tidak memiliki akses")
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

/*
Delete
*/
