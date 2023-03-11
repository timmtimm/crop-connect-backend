package users

import (
	"marketplace-backend/business/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Description string             `json:"description" bson:"description"`
	PhoneNumber string             `json:"phoneNumber" bson:"phoneNumber"`
	Password    string             `json:"password" bson:"password"`
	Role        string             `json:"role" bson:"role"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain *users.Domain) *Model {
	return &Model{
		ID:          domain.ID,
		Name:        domain.Name,
		Email:       domain.Email,
		Description: domain.Description,
		PhoneNumber: domain.PhoneNumber,
		Password:    domain.Password,
		Role:        domain.Role,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}
}

func (model *Model) ToDomain() users.Domain {
	return users.Domain{
		ID:          model.ID,
		Name:        model.Name,
		Email:       model.Email,
		Description: model.Description,
		PhoneNumber: model.PhoneNumber,
		Password:    model.Password,
		Role:        model.Role,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
