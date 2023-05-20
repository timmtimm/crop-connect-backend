package regions

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	ID          primitive.ObjectID
	Country     string // Negara
	Province    string // Provinsi
	Regency     string // Kota/Kabupaten
	District    string // Kecamatan
	Subdistrict string // Kelurahan/Desa
}

type Query struct {
	Country     string
	Province    string
	Regency     string
	District    string
	Subdistrict string
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByQuery(query Query) ([]Domain, error)
	GetProvince(country string) ([]string, error)
	GetRegency(country string, province string) ([]string, error)
	GetDistrict(country string, province string, regency string) ([]string, error)
	GetSubdistrict(country string, province string, regency string, district string) ([]Domain, error)
}

type UseCase interface {
	// Create
	Create(domain *Domain) (Domain, int, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, int, error)
	GetByQuery(query Query) ([]Domain, int, error)
	GetByCountry(country string) ([]string, int, error)
	GetByProvince(country string, province string) ([]string, int, error)
	GetByRegency(country string, province string, regency string) ([]string, int, error)
	GetByDistrict(country string, province string, regency string, district string) ([]Domain, int, error)
}
