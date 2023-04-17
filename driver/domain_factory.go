package driver

import (
	batchDomain "crop_connect/business/batchs"
	commodityDomain "crop_connect/business/commodities"
	harvestDomain "crop_connect/business/harvests"
	proposalDomain "crop_connect/business/proposals"
	regionDomain "crop_connect/business/regions"
	transactionDomain "crop_connect/business/transactions"
	treatmentRecordDomain "crop_connect/business/treatment_records"
	userDomain "crop_connect/business/users"

	batchDB "crop_connect/driver/mongo/batchs"
	commodityDB "crop_connect/driver/mongo/commodities"
	harvestDB "crop_connect/driver/mongo/harvests"
	proposalDB "crop_connect/driver/mongo/proposals"
	regionDB "crop_connect/driver/mongo/regions"
	transactionDB "crop_connect/driver/mongo/transactions"
	treatmentRecordDB "crop_connect/driver/mongo/treatment_records"
	userDB "crop_connect/driver/mongo/users"

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

func NewHarvestRepository(db *mongo.Database) harvestDomain.Repository {
	return harvestDB.NewMongoRepository(db)
}

func NewRegionRepository(db *mongo.Database) regionDomain.Repository {
	return regionDB.NewMongoRepository(db)
}
