package response

import (
	"crop_connect/business/regions"
	"crop_connect/business/users"
	regionResponse "crop_connect/controller/regions/response"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID      `json:"_id"`
	Region      regionResponse.Response `json:"region"`
	Name        string                  `json:"name"`
	Email       string                  `json:"email"`
	Description string                  `json:"description"`
	PhoneNumber string                  `json:"phoneNumber"`
	Role        string                  `json:"role"`
	CreatedAt   primitive.DateTime      `json:"createdAt"`
	UpdatedAt   primitive.DateTime      `json:"updatedAt,omitempty"`
}

func FromDomain(domain users.Domain, regionUC regions.UseCase) (User, int, error) {
	region, statusCode, err := regionUC.GetByID(domain.RegionID)
	if err != nil {
		return User{}, statusCode, err
	}

	return User{
		ID:          domain.ID,
		Region:      regionResponse.FromDomain(&region),
		Name:        domain.Name,
		Email:       domain.Email,
		Description: domain.Description,
		PhoneNumber: domain.PhoneNumber,
		Role:        domain.Role,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}, http.StatusOK, nil
}

func FromDomainArray(data []users.Domain, regionUC regions.UseCase) ([]User, int, error) {
	var response []User
	for _, domain := range data {
		user, statusCode, err := FromDomain(domain, regionUC)
		if err != nil {
			return []User{}, statusCode, err
		}

		response = append(response, user)
	}

	return response, http.StatusOK, nil
}
