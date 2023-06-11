package forgot_password

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	ID        primitive.ObjectID
	Email     string
	Token     string
	IsUsed    bool
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
	ExpiredAt primitive.DateTime
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByToken(token string) (Domain, error)
	// Update
	Update(domain *Domain) (Domain, error)
	// Delete
	HardDelete(id primitive.ObjectID) error
}

type UseCase interface {
	// Create
	Generate(appDomain string, email string) (int, error)
	// Read
	ValidateToken(token string) (int, error)
	// Update
	ResetPassword(token string, password string) (int, error)
	// Delete
}
