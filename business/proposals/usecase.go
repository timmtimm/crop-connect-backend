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
	}

	_, err = pr.proposalRepository.GetByCommodityIDAndName(domain.CommodityID, domain.Name)
	if err == nil {
		return http.StatusConflict, errors.New("nama proposal sudah digunakan")
	}

	domain.ID = primitive.NewObjectID()
	domain.IsAccepted = false
	domain.IsAvailable = true
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

/*
Delete
*/
