package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ID          primitive.ObjectID `json:"_id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Description string             `json:"description"`
	PhoneNumber string             `json:"phoneNumber"`
	Password    string             `json:"-"`
	Role        string             `json:"role"`
	CreatedAt   primitive.DateTime `json:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt"`
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByID(id primitive.ObjectID) (Domain, error)
	GetByEmail(email string) (Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
}

type UseCase interface {
	// Create
	Register(domain *Domain) (string, int, error)
	// Read
	Login(domain *Domain) (string, int, error)
	GetByID(id primitive.ObjectID) (Domain, int, error)
	// Update
	UpdateProfile(domain *Domain) (Domain, int, error)
	// Delete
}
