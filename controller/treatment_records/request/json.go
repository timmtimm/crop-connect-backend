package request

import (
	"errors"
	treatmentRecords "marketplace-backend/business/treatment_records"
	"marketplace-backend/helper"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestToFarmer struct {
	Date        string `form:"date" json:"date" validate:"required"`
	Description string `form:"description" json:"description" validate:"required"`
}

func (req *RequestToFarmer) ToDomain() (*treatmentRecords.Domain, error) {
	dateOnTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return &treatmentRecords.Domain{}, errors.New("date harus berupa tanggal")
	}

	return &treatmentRecords.Domain{
		Date:        primitive.NewDateTimeFromTime(dateOnTime),
		Description: req.Description,
	}, nil
}

func (req *RequestToFarmer) Validate() []helper.ValidationError {
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

type FillTreatmentRecord struct {
	Notes []string `form:"notes" json:"notes"`
}

type Validate struct {
	Status       string `form:"status" json:"status" validate:"required"`
	RevisionNote string `form:"revisionNote" json:"revisionNote"`
	WarningNote  string `form:"warningNote" json:"warningNote"`
}

func (req *Validate) ToDomain() *treatmentRecords.Domain {
	return &treatmentRecords.Domain{
		Status:       req.Status,
		RevisionNote: req.RevisionNote,
		WarningNote:  req.WarningNote,
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

type UpdateNotes struct {
	RevisionNote string `form:"revisionNote" json:"revisionNote"`
	WarningNote  string `form:"warningNote" json:"warningNote"`
}

func (req *UpdateNotes) ToDomain() *treatmentRecords.Domain {
	return &treatmentRecords.Domain{
		RevisionNote: req.RevisionNote,
		WarningNote:  req.WarningNote,
	}
}
