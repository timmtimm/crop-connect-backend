package request

import (
	"crop_connect/business/harvests"
	"crop_connect/helper"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmitHarvest struct {
	Date         string `form:"date" json:"date" validate:"required"`
	TotalHarvest string `form:"totalHarvest" json:"totalHarvest" validate:"required,number"`
	Condition    string `form:"condition" json:"condition" validate:"required"`
	Note1        string `form:"note1" json:"note1"`
	Note2        string `form:"note2" json:"note2"`
	Note3        string `form:"note3" json:"note3"`
	Note4        string `form:"note4" json:"note4"`
	Note5        string `form:"note5" json:"note5"`
	IsChange     string `form:"isChange" json:"isChange"`
	IsDelete     string `form:"isDelete" json:"isDelete"`
}

func (req *SubmitHarvest) ToDomain() (*harvests.Domain, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return &harvests.Domain{}, errors.New("date harus berupa tanggal")
	}

	totalHarvest, err := strconv.ParseFloat(req.TotalHarvest, 64)
	if err != nil {
		return &harvests.Domain{}, errors.New("totalHarvest harus berupa angka")
	}

	return &harvests.Domain{
		Date:         primitive.NewDateTimeFromTime(date),
		TotalHarvest: totalHarvest,
		Condition:    req.Condition,
	}, nil
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
