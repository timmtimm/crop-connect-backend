package request

import (
	"crop_connect/business/transactions"
	"crop_connect/constant"
	"crop_connect/helper"
	"crop_connect/util"
	"errors"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Create struct {
	TransactionType string `form:"transactionType" json:"transactionType" validate:"required"`
	TransactedID    string `form:"transactedID" json:"transactedID" validate:"required"`
	RegionID        string `form:"regionID" json:"regionID" validate:"required"`
	Address         string `form:"address" json:"address" validate:"required"`
}

func (req *Create) ToDomain() (*transactions.Domain, error) {
	isAvailable := util.CheckStringOnArray([]string{constant.TransactionTypeAnnuals, constant.TransactionTypePerennials}, req.TransactionType)
	if !isAvailable {
		return nil, errors.New("jenis transaksi tidak tersedia")
	}

	transactedObjID, err := primitive.ObjectIDFromHex(req.TransactedID)
	if err != nil {
		return nil, errors.New("id yang ditransaksikan tidak valid")
	}

	regionObjID, err := primitive.ObjectIDFromHex(req.RegionID)
	if err != nil {
		return nil, errors.New("id daerah tidak valid")
	}

	return &transactions.Domain{
		TransactionType: req.TransactionType,
		TransactedID:    transactedObjID,
		RegionID:        regionObjID,
		Address:         req.Address,
	}, nil
}

func (req *Create) Validate() []helper.ValidationError {
	var ve validator.ValidationErrors

	if err := validator.New().Struct(req); err != nil {
		if errors.As(err, &ve) {
			fields := structs.Fields(req)
			out := make([]helper.ValidationError, len(ve))

			for i, e := range ve {
				out[i] = helper.ValidationError{
					Field:   e.Field(),
					Message: helper.MessageForTag(e.Tag()),
				}

				out[i].Message = strings.Replace(out[i].Message, "[PARAM]", e.Param(), 1)

				for _, f := range fields {
					if f.Name() == e.Field() {
						out[i].Field = f.Tag("json")
						break
					}
				}
			}
			return out
		}
	}

	return nil
}

type Decision struct {
	Decision string `form:"decision" json:"decision"`
}

func (req *Decision) ToDomain() *transactions.Domain {
	return &transactions.Domain{
		Status: req.Decision,
	}
}
