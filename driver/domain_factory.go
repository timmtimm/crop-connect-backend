package driver

import (
	userDomain "marketplace-backend/business/users"

	userDB "marketplace-backend/driver/mongo/users"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepository(db *mongo.Database) userDomain.Repository {
	return userDB.NewMongoRepository(db)
}
