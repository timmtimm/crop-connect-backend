package commodities

import (
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommoditiesUseCase struct {
	commoditiesRepository Repository
}

func NewCommodityUseCase(cr Repository) UseCase {
	return &CommoditiesUseCase{
		commoditiesRepository: cr,
	}
}

/*
Create
*/

func (cu *CommoditiesUseCase) Create(domain *Domain) (int, error) {
	_, err := cu.commoditiesRepository.GetByNameAndFarmerID(domain.Name, domain.FarmerID)
	if err == nil {
		return http.StatusConflict, errors.New("nama komoditas telah terdaftar")
	}

	domain.ID = primitive.NewObjectID()
	domain.ImageURLs = []string{}
	domain.IsAvailable = true
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = cu.commoditiesRepository.Create(domain)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal membuat komoditas")
	}

	return http.StatusCreated, nil
}

/*
Read
*/

func (cu *CommoditiesUseCase) GetByPaginationAndQuery(query Query) ([]Domain, int, int, error) {
	commodities, totalData, err := cu.commoditiesRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, errors.New("gagal mendapatkan komoditas")
	}

	return commodities, totalData, http.StatusOK, nil
}

/*
Update
*/

func (cu *CommoditiesUseCase) Update(domain *Domain) (int, error) {
	commodity, err := cu.commoditiesRepository.GetByIDAndFarmerID(domain.ID, domain.FarmerID)
	if err != nil {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	}

	if domain.Name == commodity.Name &&
		domain.Description == commodity.Description &&
		domain.Seed == commodity.Seed &&
		domain.PlantingPeriod == commodity.PlantingPeriod &&
		domain.PricePerKg == commodity.PricePerKg {
		return http.StatusConflict, errors.New("tidak ada perubahan data")
	}

	if commodity.Name != domain.Name {
		_, err = cu.commoditiesRepository.GetByNameAndFarmerID(domain.Name, domain.FarmerID)
		if err == nil {
			return http.StatusConflict, errors.New("nama komoditas telah terdaftar")
		}
	}

	err = cu.commoditiesRepository.Delete(domain.ID)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal menghapus komoditas")
	}

	domain.ID = primitive.NewObjectID()
	domain.CreatedAt = commodity.CreatedAt
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = cu.commoditiesRepository.Create(domain)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal mengupdate komoditas")
	}

	return http.StatusOK, nil
}

/*
Delete
*/

func (cu *CommoditiesUseCase) Delete(id primitive.ObjectID, farmerID primitive.ObjectID) (int, error) {
	_, err := cu.commoditiesRepository.GetByIDAndFarmerID(id, farmerID)
	if err != nil {
		return http.StatusNotFound, errors.New("komoditas tidak ditemukan")
	}

	err = cu.commoditiesRepository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal menghapus komoditas")
	}

	return http.StatusOK, nil
}
