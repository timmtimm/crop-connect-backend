package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID          primitive.ObjectID
	RegionID    primitive.ObjectID
	Name        string
	Email       string
	Description string
	PhoneNumber string
	Password    string
	Role        string
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
}

type Query struct {
	Skip        int64
	Limit       int64
	Sort        string
	Order       int
	Name        string
	Email       string
	PhoneNumber string
	Role        string
	Province    string
	Regency     string
	District    string
	RegionID    primitive.ObjectID
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByEmail(email string) (Domain, error)
	GetByNameAndRole(name string, role string) ([]Domain, error)
	GetByQuery(query Query) ([]Domain, int, error)
	GetFarmerByID(id primitive.ObjectID) (Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
}

type UseCase interface {
	// Create
	Register(domain *Domain) (string, int, error)
	RegisterValidator(domain *Domain) (string, int, error)
	// Read
	Login(domain *Domain) (string, int, error)
	GetByID(id primitive.ObjectID) (Domain, int, error)
	GetByPaginationAndQuery(query Query) ([]Domain, int, int, error)
	GetFarmerByID(id primitive.ObjectID) (Domain, int, error)
	// Update
	UpdateProfile(domain *Domain) (Domain, int, error)
	UpdatePassword(domain *Domain, newPassword string) (Domain, int, error)
	// Delete
}
