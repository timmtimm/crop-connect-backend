package regions

import (
	"context"
	"crop_connect/business/regions"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegionRepository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) regions.Repository {
	return &RegionRepository{
		collection: db.Collection("regions"),
	}
}

/*
Create
*/

func (rr *RegionRepository) Create(domain *regions.Domain) (regions.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := rr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return regions.Domain{}, err
	}

	return *domain, err
}

/*
Read
*/

func (rr *RegionRepository) GetByID(id primitive.ObjectID) (regions.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := rr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return regions.Domain{}, err
	}

	return *result.ToDomain(), nil
}

func (rr *RegionRepository) GetByQuery(query regions.Query) ([]regions.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	filter := bson.M{}

	if query.Country != "" {
		filter["country"] = query.Country
	}

	if query.Province != "" {
		filter["province"] = query.Province
	}

	if query.Regency != "" {
		filter["regency"] = query.Regency
	}

	if query.District != "" {
		filter["district"] = query.District
	}

	if query.Subdistrict != "" {
		filter["subdistrict"] = query.Subdistrict
	}

	cursor, err := rr.collection.Find(ctx, filter)
	if err != nil {
		return []regions.Domain{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []regions.Domain{}, err
	}

	return ToDomainArray(result), nil
}

func (rr *RegionRepository) GetProvince(country string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []string

	province, err := rr.collection.Distinct(ctx, "province", bson.M{
		"country": country,
	})
	if err != nil {
		return []string{}, err
	}

	result = make([]string, len(province))
	for i, v := range province {
		result[i] = v.(string)
	}

	return result, nil
}

func (rr *RegionRepository) GetRegency(country string, province string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []string

	regency, err := rr.collection.Distinct(ctx, "regency", bson.M{
		"country":  country,
		"province": province,
	})
	if err != nil {
		return []string{}, err
	}

	result = make([]string, len(regency))
	for i, v := range regency {
		result[i] = v.(string)
	}

	return result, nil
}

func (rr *RegionRepository) GetDistrict(country string, province string, regency string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []string

	district, err := rr.collection.Distinct(ctx, "district", bson.M{
		"country":  country,
		"province": province,
		"regency":  regency,
	})
	if err != nil {
		return []string{}, err
	}

	result = make([]string, len(district))
	for i, v := range district {
		result[i] = v.(string)
	}

	return result, nil
}

func (rr *RegionRepository) GetSubdistrict(country string, province string, regency string, district string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []string

	subdistrict, err := rr.collection.Distinct(ctx, "subdistrict", bson.M{
		"country":  country,
		"province": province,
		"regency":  regency,
		"district": district,
	})
	if err != nil {
		return []string{}, err
	}

	result = make([]string, len(subdistrict))
	for i, v := range subdistrict {
		result[i] = v.(string)
	}

	return result, nil
}
