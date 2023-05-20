package forgot_password

import (
	forgotPassword "crop_connect/business/forgot_password"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	Token     string             `bson:"token"`
	IsUsed    bool               `bson:"isUsed"`
	CreatedAt primitive.DateTime `bson:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt,omitempty"`
	ExpiredAt primitive.DateTime `bson:"expiredAt"`
}

func FromDomain(domain *forgotPassword.Domain) *Model {
	return &Model{
		ID:        domain.ID,
		Email:     domain.Email,
		Token:     domain.Token,
		IsUsed:    domain.IsUsed,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
		ExpiredAt: domain.ExpiredAt,
	}
}

func (m *Model) ToDomain() forgotPassword.Domain {
	return forgotPassword.Domain{
		ID:        m.ID,
		Email:     m.Email,
		Token:     m.Token,
		IsUsed:    m.IsUsed,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		ExpiredAt: m.ExpiredAt,
	}
}
