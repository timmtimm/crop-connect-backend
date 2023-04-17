package seeds

import (
	"fmt"
	"marketplace-backend/business/regions"
	_mongo "marketplace-backend/driver/mongo"

	"go.mongodb.org/mongo-driver/mongo"
)

func SeedDatabase(db *mongo.Database, regionUC regions.UseCase) {
	if check, _ := _mongo.CheckCollectionExist(db, "regions"); !check {
		fmt.Println("Seeding regions...")
		SeedRegion(regionUC)
	}
}
