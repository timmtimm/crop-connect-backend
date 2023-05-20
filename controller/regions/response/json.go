package response

import (
	"crop_connect/business/regions"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	ID          primitive.ObjectID `json:"_id"`
	Country     string             `json:"country"`     // Negara
	Province    string             `json:"province"`    // Provinsi
	Regency     string             `json:"regency"`     // Kabupaten
	District    string             `json:"district"`    // Kecamatan
	Subdistrict string             `json:"subdistrict"` // Kelurahan
}

func FromDomain(domain *regions.Domain) Response {
	return Response{
		ID:          domain.ID,
		Country:     domain.Country,
		Province:    domain.Province,
		Regency:     domain.Regency,
		District:    domain.District,
		Subdistrict: domain.Subdistrict,
	}
}

func FromDomainArray(domain []regions.Domain) []Response {
	var response []Response
	for _, value := range domain {
		response = append(response, FromDomain(&value))
	}
	return response
}
