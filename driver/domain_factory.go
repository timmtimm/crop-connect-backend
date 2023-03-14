package driver

import (
	commodityDomain "marketplace-backend/business/commodities"
	proposalDomain "marketplace-backend/business/proposals"
	userDomain "marketplace-backend/business/users"

	commodityDB "marketplace-backend/driver/mongo/commodities"
	proposalDB "marketplace-backend/driver/mongo/proposals"
	userDB "marketplace-backend/driver/mongo/users"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepository(db *mongo.Database) userDomain.Repository {
	return userDB.NewMongoRepository(db)
}

func NewCommodityRepository(db *mongo.Database) commodityDomain.Repository {
	return commodityDB.NewMongoRepository(db)
}

func NewProposalRepository(db *mongo.Database) proposalDomain.Repository {
	return proposalDB.NewMongoRepository(db)
}
