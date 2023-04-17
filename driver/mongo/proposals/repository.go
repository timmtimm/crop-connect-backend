package proposals

import (
	"context"
	"marketplace-backend/business/proposals"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProposalRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) proposals.Repository {
	return &ProposalRepository{
		collection: db.Collection("proposals"),
	}
}

/*
Create
*/

func (pr *ProposalRepository) Create(domain *proposals.Domain) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := pr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return proposals.Domain{}, err
	}

	return *domain, err
}

/*
Read
*/

func (pr *ProposalRepository) GetByID(id primitive.ObjectID) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := pr.collection.FindOne(ctx, bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

func (pr *ProposalRepository) GetByIDWithoutDeleted(id primitive.ObjectID) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := pr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (pr *ProposalRepository) GetByCommodityID(commodityID primitive.ObjectID) ([]proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := pr.collection.Find(ctx, bson.M{
		"commodityID": commodityID,
		"deletedAt":   bson.M{"$exists": false},
	})
	if err != nil {
		return []proposals.Domain{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []proposals.Domain{}, err
	}

	return ToDomainArray(result), err
}

func (pr *ProposalRepository) GetByCommodityIDAndAvailability(commodityID primitive.ObjectID, status string) ([]proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := pr.collection.Find(ctx, bson.M{
		"commodityID": commodityID,
		"status":      status,
		"deletedAt":   bson.M{"$exists": false},
	})
	if err != nil {
		return []proposals.Domain{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []proposals.Domain{}, err
	}

	return ToDomainArray(result), err
}

func (pr *ProposalRepository) GetByCommodityIDAndName(commodityID primitive.ObjectID, name string) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := pr.collection.FindOne(ctx, bson.M{
		"commodityID": commodityID,
		"name":        name,
		"deletedAt":   bson.M{"$exists": false},
	}).Decode(&result)

	return result.ToDomain(), err
}

/*
Update
*/

func (pr *ProposalRepository) Update(domain *proposals.Domain) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := pr.collection.UpdateOne(ctx, bson.M{
		"_id":       domain.ID,
		"deletedAt": bson.M{"$exists": false},
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return proposals.Domain{}, err
	}

	return *domain, nil
}

func (pr *ProposalRepository) UnsetRejectReason(id primitive.ObjectID) (proposals.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := pr.collection.UpdateOne(ctx, bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}, bson.M{
		"$unset": bson.M{"rejectReason": ""},
	})
	if err != nil {
		return proposals.Domain{}, err
	}

	updatedProposal, err := pr.GetByID(id)
	if err != nil {
		return proposals.Domain{}, err
	}

	return updatedProposal, nil
}

/*
Delete
*/

func (pr *ProposalRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := pr.collection.UpdateOne(ctx, bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}, bson.M{
		"$set": bson.M{
			"deletedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (pr *ProposalRepository) DeleteByCommodityID(commodityID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := pr.collection.UpdateMany(ctx, bson.M{
		"commodityID": commodityID,
		"deletedAt":   bson.M{"$exists": false},
	}, bson.M{
		"$set": bson.M{
			"deletedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	})

	return err
}
