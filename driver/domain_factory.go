package driver

import (
	batchDomain "marketplace-backend/business/batchs"
	commodityDomain "marketplace-backend/business/commodities"
	proposalDomain "marketplace-backend/business/proposals"
	transactionDomain "marketplace-backend/business/transactions"
	treatmentRecordDomain "marketplace-backend/business/treatment_records"
	userDomain "marketplace-backend/business/users"

	batchDB "marketplace-backend/driver/mongo/batchs"
	commodityDB "marketplace-backend/driver/mongo/commodities"
	proposalDB "marketplace-backend/driver/mongo/proposals"
	transactionDB "marketplace-backend/driver/mongo/transactions"
	treatmentRecordDB "marketplace-backend/driver/mongo/treatment_records"
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

func NewTransactionRepository(db *mongo.Database) transactionDomain.Repository {
	return transactionDB.NewMongoRepository(db)
}

func NewBatchRepository(db *mongo.Database) batchDomain.Repository {
	return batchDB.NewMongoRepository(db)
}

func NewTreatmentRecordRepository(db *mongo.Database) treatmentRecordDomain.Repository {
	return treatmentRecordDB.NewMongoRepository(db)
}
