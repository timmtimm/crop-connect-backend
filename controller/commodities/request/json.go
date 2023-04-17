package request

import (
	"errors"
	"marketplace-backend/business/commodities"
	"marketplace-backend/helper"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

type Create struct {
	Name           string `form:"name" json:"name" validate:"required,min=3,max=100"`
	Description    string `form:"description" json:"description"`
	Seed           string `form:"seed" json:"seed" validate:"required,min=3,max=100"`
	PlantingPeriod int    `form:"plantingPeriod" json:"plantingPeriod" validate:"required,number"`
	PricePerKg     int    `form:"pricePerKg" json:"pricePerKg" validate:"required,number"`
}

func (req *Create) ToDomain() *commodities.Domain {
	return &commodities.Domain{
		Name:           req.Name,
		Description:    req.Description,
		Seed:           req.Seed,
		PlantingPeriod: req.PlantingPeriod,
		PricePerKg:     req.PricePerKg,
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

type Update struct {
	Name           string `form:"name" json:"name" validate:"required,min=3,max=100"`
	Description    string `form:"description" json:"description"`
	Seed           string `form:"seed" json:"seed" validate:"required,min=3,max=100"`
	PlantingPeriod int    `form:"plantingPeriod" json:"plantingPeriod" validate:"required,number"`
	PricePerKg     int    `form:"pricePerKg" json:"pricePerKg" validate:"required,number"`
	IsChange       []bool `json:"isChange" form:"isChange"`
	IsDelete       []bool `json:"isDelete" form:"isDelete"`
}

func (req *Update) ToDomain() *commodities.Domain {
	return &commodities.Domain{
		Name:           req.Name,
		Description:    req.Description,
		Seed:           req.Seed,
		PlantingPeriod: req.PlantingPeriod,
		PricePerKg:     req.PricePerKg,
	}
}

func (req *Update) Validate() []helper.ValidationError {
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
