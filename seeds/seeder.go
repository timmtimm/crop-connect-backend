package seeds

import (
	"crop_connect/business/regions"
	_mongo "crop_connect/driver/mongo"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func SeedDatabase(db *mongo.Database, regionUC regions.UseCase) {
	if check, _ := _mongo.CheckCollectionExist(db, "regions"); !check {
		fmt.Println("Seeding regions...")
		SeedRegion(regionUC)
	}
}
