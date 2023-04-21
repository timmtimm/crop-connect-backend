package transactions

import (
	"crop_connect/business/transactions"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID         primitive.ObjectID `bson:"_id"`
	BuyerID    primitive.ObjectID `bson:"buyerID"`
	ProposalID primitive.ObjectID `bson:"proposalID"`
	RegionID   primitive.ObjectID `bson:"regionID"`
	Address    string             `bson:"address"`
	Status     string             `bson:"status"`
	TotalPrice float64            `bson:"totalPrice"`
	CreatedAt  primitive.DateTime `bson:"createdAt"`
	UpdatedAt  primitive.DateTime `bson:"updatedAt,omitempty"`
}

func FromDomain(domain *transactions.Domain) *Model {
	return &Model{
		ID:         domain.ID,
		BuyerID:    domain.BuyerID,
		ProposalID: domain.ProposalID,
		RegionID:   domain.RegionID,
		Address:    domain.Address,
		Status:     domain.Status,
		TotalPrice: domain.TotalPrice,
		CreatedAt:  domain.CreatedAt,
		UpdatedAt:  domain.UpdatedAt,
	}
}

func (model *Model) ToDomain() transactions.Domain {
	return transactions.Domain{
		ID:         model.ID,
		BuyerID:    model.BuyerID,
		ProposalID: model.ProposalID,
		RegionID:   model.RegionID,
		Address:    model.Address,
		Status:     model.Status,
		TotalPrice: model.TotalPrice,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
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
	Month            int `bson:"month"`
	TotalAccepted    int `bson:"totalAccepted"`
	TotalTransaction int `bson:"totalTransaction"`
	TotalIncome      int `bson:"totalIncome"`
}

func (model *StatisticModel) ToStatistic() transactions.Statistic {
	return transactions.Statistic{
		Month:            model.Month,
		TotalAccepted:    model.TotalAccepted,
		TotalTransaction: model.TotalTransaction,
		TotalIncome:      model.TotalIncome,
	}
}

func ToStatisticArray(models []StatisticModel) []transactions.Statistic {
	var domains []transactions.Statistic
	for _, model := range models {
		domains = append(domains, model.ToStatistic())
	}
	return domains
}
