package seeds

import (
	regionDomain "crop_connect/business/regions"
	_mongo "crop_connect/driver/mongo"
	regionSeed "crop_connect/seeds/regions"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func SeedDatabase(db *mongo.Database, regionUC regionDomain.UseCase) {
	if check, _ := _mongo.CheckCollectionExist(db, "regions"); !check {
		fmt.Println("Seeding regions...")
		regionSeed.Seed(regionUC)
	}
}
