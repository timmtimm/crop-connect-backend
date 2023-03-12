package driver

import (
	commodityDomain "marketplace-backend/business/commodities"
	userDomain "marketplace-backend/business/users"

	commodityDB "marketplace-backend/driver/mongo/commodities"
	userDB "marketplace-backend/driver/mongo/users"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepository(db *mongo.Database) userDomain.Repository {
	return userDB.NewMongoRepository(db)
}

func NewCommodityRepository(db *mongo.Database) commodityDomain.Repository {
	return commodityDB.NewMongoRepository(db)
}
