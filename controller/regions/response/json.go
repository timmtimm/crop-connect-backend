package response

import (
	"crop_connect/business/regions"
)

type Response struct {
	Country     string `json:"country"`               //negara
	Province    string `json:"province"`              // provinsi
	Regency     string `json:"regency,omitempty"`     // kabupaten
	District    string `json:"district,omitempty"`    // kecamatan
	Subdistrict string `json:"subdistrict,omitempty"` // kelurahan
}

func FromDomain(domain *regions.Domain) Response {
	return Response{
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
