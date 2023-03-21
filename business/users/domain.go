package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID          primitive.ObjectID
	Name        string
	Email       string
	Description string
	PhoneNumber string
	Password    string
	Role        string
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByEmail(email string) (Domain, error)
	GetByNameAndRole(name string, role string) ([]Domain, error)
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
	GetFarmerByName(name string) ([]Domain, int, error)
	// Update
	UpdateProfile(domain *Domain) (Domain, int, error)
	// Delete
}
