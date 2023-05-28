package request

import (
	"crop_connect/business/proposals"
	"crop_connect/helper"
	"errors"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Create struct {
	RegionID              string  `form:"regionID" json:"regionID" validate:"required"`
	Name                  string  `form:"name" json:"name" validate:"required"`
	Description           string  `form:"description" json:"description"`
	EstimatedTotalHarvest float64 `form:"estimatedTotalHarvest" json:"estimatedTotalHarvest" validate:"required,number"`
	PlantingArea          float64 `form:"plantingArea" json:"plantingArea" validate:"required,number"`
	Address               string  `form:"address" json:"address" validate:"required"`
	IsAvailable           bool    `form:"isAvailable" json:"isAvailable"`
}

func (req *Create) ToDomain() (*proposals.Domain, error) {
	regionObjID, err := primitive.ObjectIDFromHex(req.RegionID)
	if err != nil {
		return nil, errors.New("id daerah tidak valid")
	}

	return &proposals.Domain{
		RegionID:              regionObjID,
		Name:                  req.Name,
		Description:           req.Description,
		EstimatedTotalHarvest: req.EstimatedTotalHarvest,
		PlantingArea:          req.PlantingArea,
		Address:               req.Address,
		IsAvailable:           req.IsAvailable,
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

type Validate struct {
	Status       string `form:"status" json:"status"`
	RejectReason string `form:"rejectReason" json:"rejectReason"`
}

func (req *Validate) ToDomain() *proposals.Domain {
	return &proposals.Domain{
		Status:       req.Status,
		RejectReason: req.RejectReason,
	}
}
