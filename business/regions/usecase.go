package regions

import (
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegionUseCase struct {
	regionRepository Repository
}

func NewRegionUseCase(rr Repository) UseCase {
	return &RegionUseCase{
		regionRepository: rr,
	}
}

/*
Create
*/

func (ru *RegionUseCase) Create(domain *Domain) (Domain, int, error) {
	domain.ID = primitive.NewObjectID()

	region, err := ru.regionRepository.Create(domain)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal membuat daerah")
	}

	return region, http.StatusCreated, nil
}

/*
Read
*/

func (ru *RegionUseCase) GetByID(id primitive.ObjectID) (Domain, int, error) {
	region, err := ru.regionRepository.GetByID(id)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("daerah tidak ditemukan")
	}

	return region, http.StatusOK, nil
}

func (ru *RegionUseCase) GetByQuery(query Query) ([]Domain, int, error) {
	regions, err := ru.regionRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, http.StatusInternalServerError, err
	}

	return regions, http.StatusOK, nil
}

func (ru *RegionUseCase) GetByCountry(country string) ([]string, int, error) {
	regions, err := ru.regionRepository.GetProvince(country)
	if err != nil {
		return []string{}, http.StatusInternalServerError, err
	}

	return regions, http.StatusOK, nil
}

func (ru *RegionUseCase) GetByProvince(country string, province string) ([]string, int, error) {
	regions, err := ru.regionRepository.GetRegency(country, province)
	if err != nil {
		return []string{}, http.StatusInternalServerError, err
	}

	return regions, http.StatusOK, nil
}

func (ru *RegionUseCase) GetByRegency(country string, province string, regency string) ([]string, int, error) {
	regions, err := ru.regionRepository.GetDistrict(country, province, regency)
	if err != nil {
		return []string{}, http.StatusInternalServerError, err
	}

	return regions, http.StatusOK, nil
}

func (ru *RegionUseCase) GetByDistrict(country string, province string, regency string, district string) ([]string, int, error) {
	regions, err := ru.regionRepository.GetSubdistrict(country, province, regency, district)
	if err != nil {
		return []string{}, http.StatusInternalServerError, err
	}

	return regions, http.StatusOK, nil
}
