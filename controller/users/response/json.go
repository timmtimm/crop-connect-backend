package response

import (
	"marketplace-backend/business/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Description string             `json:"description"`
	PhoneNumber string             `json:"phoneNumber"`
	Role        string             `json:"role"`
	CreatedAt   primitive.DateTime `json:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt,omitempty"`
}

func FromDomain(domain users.Domain) User {
	return User{
		Name:        domain.Name,
		Email:       domain.Email,
		Description: domain.Description,
		PhoneNumber: domain.PhoneNumber,
		Role:        domain.Role,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}

func FromDomainArray(data []users.Domain) []User {
	var array []User
	for _, v := range data {
		array = append(array, FromDomain(v))
	}
	return array
}
