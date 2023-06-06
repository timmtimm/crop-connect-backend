package response

import (
	"crop_connect/business/batchs"
	"crop_connect/business/commodities"
	"crop_connect/business/proposals"
	"crop_connect/business/regions"
	"crop_connect/business/transactions"
	"crop_connect/business/users"
	"crop_connect/constant"
	"net/http"

	batchResponse "crop_connect/controller/batchs/response"
	commodityResponse "crop_connect/controller/commodities/response"
	proposalResponse "crop_connect/controller/proposals/response"
	regionResponse "crop_connect/controller/regions/response"
	userResponse "crop_connect/controller/users/response"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Buyer struct {
	ID              primitive.ObjectID                 `json:"_id"`
	Region          regionResponse.Response            `json:"region"`
	Commodity       commodityResponse.Commodity        `json:"commodity"`
	Proposal        proposalResponse.Buyer             `json:"proposal"`
	Batch           batchResponse.BatchWithoutProposal `json:"batch"`
	Address         string                             `json:"address"`
	TransactionType string                             `json:"transactionType"`
	Status          string                             `json:"status"`
	TotalPrice      float64                            `json:"totalPrice"`
	CreatedAt       primitive.DateTime                 `json:"createdAt"`
}

func FromDomainToBuyer(domain *transactions.Domain, batchUC batchs.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) (Buyer, int, error) {
	batch, statusCode, err := batchUC.GetByID(domain.BatchID)
	if err != nil {
		return Buyer{}, statusCode, err
	}

	proposal, statusCode, err := proposalUC.GetByIDWithoutDeleted(batch.ProposalID)
	if err != nil {
		return Buyer{}, statusCode, err
	}

	commodityDomain, statusCode, err := commodityUC.GetByIDWithoutDeleted(proposal.CommodityID)
	if err != nil {
		return Buyer{}, statusCode, err
	}

	commodity, statusCode, err := commodityResponse.FromDomain(commodityDomain, userUC, regionUC)
	if err != nil {
		return Buyer{}, statusCode, err
	}

	region, statusCode, err := regionUC.GetByID(domain.RegionID)
	if err != nil {
		return Buyer{}, statusCode, err
	}

	return Buyer{
		ID:              domain.ID,
		Region:          regionResponse.FromDomain(&region),
		Commodity:       commodity,
		Proposal:        proposalResponse.FromDomainToBuyer(&proposal),
		Batch:           batchResponse.FromDomainWithoutProposal(&batch),
		Address:         domain.Address,
		Status:          domain.Status,
		TransactionType: domain.TransactionType,
		TotalPrice:      domain.TotalPrice,
		CreatedAt:       domain.CreatedAt,
	}, http.StatusOK, nil
}

func FromDomainArrayToBuyer(domain []transactions.Domain, batchUC batchs.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) ([]Buyer, int, error) {
	var buyers []Buyer
	for _, value := range domain {
		buyer, statusCode, err := FromDomainToBuyer(&value, batchUC, proposalUC, commodityUC, userUC, regionUC)
		if err != nil {
			return []Buyer{}, statusCode, err
		}

		buyers = append(buyers, buyer)
	}

	return buyers, http.StatusOK, nil
}

type All struct {
	ID              primitive.ObjectID          `json:"_id"`
	Region          regionResponse.Response     `json:"region"`
	Buyer           userResponse.User           `json:"buyer"`
	Commodity       commodityResponse.Commodity `json:"commodity"`
	Proposal        proposalResponse.Buyer      `json:"proposal"`
	Batch           batchResponse.Batch         `json:"batch"`
	Address         string                      `json:"address"`
	Status          string                      `json:"status"`
	TransactionType string                      `json:"transactionType"`
	TotalPrice      float64                     `json:"totalPrice"`
	CreatedAt       primitive.DateTime          `json:"createdAt"`
}

func FromDomainToFarmer(domain *transactions.Domain, batchUC batchs.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) (All, int, error) {
	batch, statusCode, err := batchUC.GetByID(domain.BatchID)
	if err != nil {
		return All{}, statusCode, err
	}

	proposal, statusCode, err := proposalUC.GetByIDWithoutDeleted(batch.ProposalID)
	if err != nil {
		return All{}, statusCode, err
	}

	commodityDomain, statusCode, err := commodityUC.GetByIDWithoutDeleted(proposal.CommodityID)
	if err != nil {
		return All{}, statusCode, err
	}

	commodity, statusCode, err := commodityResponse.FromDomain(commodityDomain, userUC, regionUC)
	if err != nil {
		return All{}, statusCode, err
	}

	buyer, statusCode, err := userUC.GetByID(domain.BuyerID)
	if err != nil {
		return All{}, statusCode, err
	}

	buyerResponse, statusCode, err := userResponse.FromDomain(buyer, regionUC)
	if err != nil {
		return All{}, statusCode, err
	}

	region, statusCode, err := regionUC.GetByID(domain.RegionID)
	if err != nil {
		return All{}, statusCode, err
	}

	return All{
		ID:              domain.ID,
		Region:          regionResponse.FromDomain(&region),
		Buyer:           buyerResponse,
		Commodity:       commodity,
		Proposal:        proposalResponse.FromDomainToBuyer(&proposal),
		Address:         domain.Address,
		TransactionType: domain.TransactionType,
		Status:          domain.Status,
		TotalPrice:      domain.TotalPrice,
		CreatedAt:       domain.CreatedAt,
	}, http.StatusOK, nil
}

func FromDomainArrayToFarmer(domains []transactions.Domain, batchUC batchs.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) ([]All, int, error) {
	var all []All
	for _, value := range domains {
		allResponse, statusCode, err := FromDomainToFarmer(&value, batchUC, proposalUC, commodityUC, userUC, regionUC)
		if err != nil {
			return []All{}, statusCode, err
		}

		all = append(all, allResponse)
	}

	return all, http.StatusOK, nil
}

type Statistic struct {
	Month            int     `json:"month"`
	TotalAccepted    int     `json:"totalAccepted"`
	TotalTransaction int     `json:"totalTransaction"`
	TotalIncome      float64 `json:"totalIncome"`
	TotalWeight      float64 `json:"totalWeight"`
	TotalUniqueBuyer int     `json:"totalUniqueBuyer"`
}

func FromDomainArrayToStatistic(domain []transactions.Statistic) []Statistic {
	var statistics []Statistic
	for _, value := range domain {
		statistics = append(statistics, Statistic{
			Month:            value.Month,
			TotalAccepted:    value.TotalAccepted,
			TotalTransaction: value.TotalTransaction,
			TotalIncome:      value.TotalIncome,
			TotalWeight:      value.TotalWeight,
			TotalUniqueBuyer: value.TotalUniqueBuyer,
		})
	}

	return statistics
}

type TotalTransactionByProvince struct {
	Province         string `json:"province"`
	TotalAccepted    int    `json:"totalAccepted"`
	TotalTransaction int    `json:"totalTransaction"`
}

func FromDomainArrayToStatisticProvince(domain []transactions.TotalTransactionByProvince) []TotalTransactionByProvince {
	var totalTransactionByProvinces []TotalTransactionByProvince
	for _, value := range domain {
		totalTransactionByProvinces = append(totalTransactionByProvinces, TotalTransactionByProvince{
			Province:         value.Province,
			TotalAccepted:    value.TotalAccepted,
			TotalTransaction: value.TotalTransaction,
		})
	}

	return totalTransactionByProvinces
}

type StatisticTopCommodity struct {
	Commodity commodityResponse.Commodity `json:"commodity"`
	Total     int                         `json:"total"`
}

func FromDomainArrayToStatisticTopCommodity(domain []transactions.StatisticTopCommodity, userUC users.UseCase, regionUC regions.UseCase) ([]StatisticTopCommodity, int, error) {
	var statistics []StatisticTopCommodity
	for _, value := range domain {
		commodity, statusCode, err := commodityResponse.FromDomain(value.Commodity, userUC, regionUC)
		if err != nil {
			return []StatisticTopCommodity{}, statusCode, err
		}

		statistics = append(statistics, StatisticTopCommodity{
			Commodity: commodity,
			Total:     value.Total,
		})
	}

	return statistics, http.StatusOK, nil
}

type TransactionStatisticForCommodityPage struct {
	TotalTransaction int     `json:"totalTransaction"`
	TotalWeight      float64 `json:"totalWeight"`
}

func FromDomainToTransactionStatisticForCommodityPage(totalTransaction int, totalWeight float64) TransactionStatisticForCommodityPage {
	return TransactionStatisticForCommodityPage{
		TotalTransaction: totalTransaction,
		TotalWeight:      totalWeight,
	}
}

type TransactionAnnuals struct {
	ID              primitive.ObjectID                 `json:"_id"`
	Region          regionResponse.Response            `json:"region"`
	Buyer           userResponse.User                  `json:"buyer"`
	Commodity       commodityResponse.Commodity        `json:"commodity"`
	Proposal        proposalResponse.Buyer             `json:"proposal"`
	Batch           batchResponse.BatchWithoutProposal `json:"batch"`
	Address         string                             `json:"address"`
	TransactionType string                             `json:"transactionType"`
	Status          string                             `json:"status"`
	TotalPrice      float64                            `json:"totalPrice"`
	CreatedAt       primitive.DateTime                 `json:"createdAt"`
}

func ConvertToTransactionResponse(domain *transactions.Domain, batchUC batchs.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) (interface{}, int, error) {
	buyer, statusCode, err := userUC.GetByID(domain.BuyerID)
	if err != nil {
		return TransactionAnnuals{}, statusCode, err
	}

	buyerResponse, statusCode, err := userResponse.FromDomain(buyer, regionUC)
	if err != nil {
		return TransactionAnnuals{}, statusCode, err
	}

	region, statusCode, err := regionUC.GetByID(domain.RegionID)
	if err != nil {
		return TransactionAnnuals{}, statusCode, err
	}

	response := TransactionAnnuals{
		ID:              domain.ID,
		Region:          regionResponse.FromDomain(&region),
		Buyer:           buyerResponse,
		Address:         domain.Address,
		TransactionType: domain.TransactionType,
		Status:          domain.Status,
		TotalPrice:      domain.TotalPrice,
		CreatedAt:       domain.CreatedAt,
	}

	if domain.TransactionType == constant.TransactionTypeAnnuals {
		proposal, statusCode, err := proposalUC.GetByIDWithoutDeleted(domain.ProposalID)
		if err != nil {
			return TransactionAnnuals{}, statusCode, err
		}

		commodityDomain, statusCode, err := commodityUC.GetByIDWithoutDeleted(proposal.CommodityID)
		if err != nil {
			return TransactionAnnuals{}, statusCode, err
		}

		commodityForResponse, statusCode, err := commodityResponse.FromDomain(commodityDomain, userUC, regionUC)
		if err != nil {
			return TransactionAnnuals{}, statusCode, err
		}

		response.Commodity = commodityForResponse
		response.Proposal = proposalResponse.FromDomainToBuyer(&proposal)

		if domain.Status == constant.TransactionStatusAccepted {
			batch, statusCode, err := batchUC.GetByID(domain.BatchID)
			if err != nil {
				return TransactionAnnuals{}, statusCode, err
			}

			response.Batch = batchResponse.FromDomainWithoutProposal(&batch)
		}
	} else if domain.TransactionType == constant.TransactionTypePerennials {
		batch, statusCode, err := batchUC.GetByID(domain.BatchID)
		if err != nil {
			return TransactionAnnuals{}, statusCode, err
		}

		proposal, statusCode, err := proposalUC.GetByIDWithoutDeleted(domain.ProposalID)
		if err != nil {
			return TransactionAnnuals{}, statusCode, err
		}

		commodityDomain, statusCode, err := commodityUC.GetByIDWithoutDeleted(proposal.CommodityID)
		if err != nil {
			return TransactionAnnuals{}, statusCode, err
		}

		commodityForResponse, statusCode, err := commodityResponse.FromDomain(commodityDomain, userUC, regionUC)
		if err != nil {
			return TransactionAnnuals{}, statusCode, err
		}

		response.Commodity = commodityForResponse
		response.Proposal = proposalResponse.FromDomainToBuyer(&proposal)
		response.Batch = batchResponse.FromDomainWithoutProposal(&batch)
	}

	return response, http.StatusOK, nil
}

func FromArrayToResponseArray(domain []transactions.Domain, batchUC batchs.UseCase, proposalUC proposals.UseCase, commodityUC commodities.UseCase, userUC users.UseCase, regionUC regions.UseCase) ([]interface{}, int, error) {
	var response []interface{}
	for _, value := range domain {
		transactionResponse, statusCode, err := ConvertToTransactionResponse(&value, batchUC, proposalUC, commodityUC, userUC, regionUC)
		if err != nil {
			return []interface{}{}, statusCode, err
		}

		response = append(response, transactionResponse)
	}

	return response, http.StatusOK, nil
}
