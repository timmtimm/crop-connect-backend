package transactions

import (
	"crop_connect/business/transactions"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID              primitive.ObjectID `bson:"_id"`
	BuyerID         primitive.ObjectID `bson:"buyerID"`
	RegionID        primitive.ObjectID `bson:"regionID"`
	TransactionType string             `bson:"transactionType"`
	TransactedID    primitive.ObjectID `bson:"transactedID"`
	Address         string             `bson:"address"`
	Status          string             `bson:"status"`
	TotalPrice      float64            `bson:"totalPrice"`
	CreatedAt       primitive.DateTime `bson:"createdAt"`
	UpdatedAt       primitive.DateTime `bson:"updatedAt,omitempty"`
}

func FromDomain(domain *transactions.Domain) *Model {
	return &Model{
		ID:              domain.ID,
		BuyerID:         domain.BuyerID,
		TransactionType: domain.TransactionType,
		TransactedID:    domain.TransactedID,
		RegionID:        domain.RegionID,
		Address:         domain.Address,
		Status:          domain.Status,
		TotalPrice:      domain.TotalPrice,
		CreatedAt:       domain.CreatedAt,
		UpdatedAt:       domain.UpdatedAt,
	}
}

func (model *Model) ToDomain() transactions.Domain {
	return transactions.Domain{
		ID:              model.ID,
		BuyerID:         model.BuyerID,
		TransactionType: model.TransactionType,
		TransactedID:    model.TransactedID,
		RegionID:        model.RegionID,
		Address:         model.Address,
		Status:          model.Status,
		TotalPrice:      model.TotalPrice,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func ToDomainArray(models []Model) []transactions.Domain {
	var domains []transactions.Domain
	for _, model := range models {
		domains = append(domains, model.ToDomain())
	}
	return domains
}

type StatisticModel struct {
	Month            int     `bson:"month"`
	TotalAccepted    int     `bson:"totalAccepted"`
	TotalTransaction int     `bson:"totalTransaction"`
	TotalIncome      float64 `bson:"totalIncome"`
	TotalWeight      float64 `bson:"totalWeight"`
	TotalUniqueBuyer int     `bson:"totalUniqueBuyer"`
}

func (model *StatisticModel) ToStatistic() transactions.Statistic {
	return transactions.Statistic{
		Month:            model.Month,
		TotalAccepted:    model.TotalAccepted,
		TotalTransaction: model.TotalTransaction,
		TotalIncome:      model.TotalIncome,
		TotalWeight:      model.TotalWeight,
		TotalUniqueBuyer: model.TotalUniqueBuyer,
	}
}

func ToStatisticArray(models []StatisticModel) []transactions.Statistic {
	var domains []transactions.Statistic
	for _, model := range models {
		domains = append(domains, model.ToStatistic())
	}
	return domains
}

type TotalTransactionByProvince struct {
	Province         string `bson:"_id"`
	TotalAccepted    int    `bson:"totalAccepted"`
	TotalTransaction int    `bson:"totalTransaction"`
}

func (model *TotalTransactionByProvince) ToTotalTransactionByProvince() transactions.TotalTransactionByProvince {
	return transactions.TotalTransactionByProvince{
		Province:         model.Province,
		TotalAccepted:    model.TotalAccepted,
		TotalTransaction: model.TotalTransaction,
	}
}

func ToTotalTransactionByProvinceArray(models []TotalTransactionByProvince) []transactions.TotalTransactionByProvince {
	var domains []transactions.TotalTransactionByProvince
	for _, model := range models {
		domains = append(domains, model.ToTotalTransactionByProvince())
	}
	return domains
}

type TotalTransactionWithWeight struct {
	TotalTransaction int     `bson:"totalTransaction"`
	TotalWeight      float64 `bson:"totalWeight"`
}
