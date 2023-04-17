package regions

import (
	"crop_connect/business/regions"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID          primitive.ObjectID `bson:"_id"`
	Country     string             `bson:"country"`     //negara
	Province    string             `bson:"province"`    // provinsi
	Regency     string             `bson:"regency"`     // kabupaten
	District    string             `bson:"district"`    // kecamatan
	Subdistrict string             `bson:"subdistrict"` // kelurahan
}

func FromDomain(domain *regions.Domain) *Model {
	return &Model{
		ID:          domain.ID,
		Country:     domain.Country,
		Province:    domain.Province,
		Regency:     domain.Regency,
		District:    domain.District,
		Subdistrict: domain.Subdistrict,
	}
}

func (m *Model) ToDomain() *regions.Domain {
	return &regions.Domain{
		ID:          m.ID,
		Country:     m.Country,
		Province:    m.Province,
		Regency:     m.Regency,
		District:    m.District,
		Subdistrict: m.Subdistrict,
	}
}

func ToDomainArray(model []Model) []regions.Domain {
	var domain []regions.Domain
	for _, v := range model {
		domain = append(domain, *v.ToDomain())
	}
	return domain
}

func InterfaceToDomain(data interface{}) (*regions.Domain, error) {
	model, ok := data.(Model)
	if !ok {
		return nil, errors.New("error when casting interface to model")
	}

	return model.ToDomain(), nil
}

func InterfaceToDomainArray(data []interface{}) ([]regions.Domain, error) {
	result := []regions.Domain{}

	for _, v := range data {
		domain, err := InterfaceToDomain(v)
		if err != nil {
			return nil, err
		}

		result = append(result, *domain)
	}

	return result, nil
}
