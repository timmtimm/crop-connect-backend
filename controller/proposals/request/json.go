package request

import (
	"errors"
	"marketplace-backend/business/proposals"
	"marketplace-backend/helper"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

type Create struct {
	Name                  string  `form:"name" json:"name" validate:"required"`
	EstimatedTotalHarvest float64 `form:"estimatedTotalHarvest" json:"estimatedTotalHarvest" validate:"required,number"`
	PlantingArea          float64 `form:"plantingArea" json:"plantingArea" validate:"required,number"`
	Address               string  `form:"address" json:"address" validate:"required"`
}

func (req *Create) ToDomain() *proposals.Domain {
	return &proposals.Domain{
		Name:                  req.Name,
		EstimatedTotalHarvest: req.EstimatedTotalHarvest,
		PlantingArea:          req.PlantingArea,
		Address:               req.Address,
	}
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
