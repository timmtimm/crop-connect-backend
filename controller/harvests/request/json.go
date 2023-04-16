package request

import (
	"errors"
	"marketplace-backend/business/harvests"
	"marketplace-backend/helper"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmitHarvest struct {
	Date         time.Time `form:"date" json:"date" validate:"required"`
	TotalHarvest float64   `form:"totalHarvest" json:"totalHarvest" validate:"required,number"`
	Condition    string    `form:"condition" json:"condition" validate:"required"`
	Notes        []string  `form:"notes" json:"notes" validate:"required"`
}

func (req *SubmitHarvest) ToDomain() *harvests.Domain {
	return &harvests.Domain{
		Date:         primitive.NewDateTimeFromTime(req.Date),
		TotalHarvest: req.TotalHarvest,
		Condition:    req.Condition,
	}
}

func (req *SubmitHarvest) Validate() []helper.ValidationError {
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
	Status       string `form:"status" json:"status" validate:"required"`
	RevisionNote string `form:"revisionNote" json:"revisionNote"`
}

func (req *Validate) ToDomain() *harvests.Domain {
	return &harvests.Domain{
		Status:       req.Status,
		RevisionNote: req.RevisionNote,
	}
}

func (req *Validate) Validate() []helper.ValidationError {
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
