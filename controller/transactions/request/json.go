package request

import (
	"errors"
	"marketplace-backend/business/transactions"
	"marketplace-backend/helper"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

type Create struct {
	Address string `form:"address" json:"address" validate:"required"`
}

func (req *Create) ToDomain() *transactions.Domain {
	return &transactions.Domain{
		Address: req.Address,
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

type Decision struct {
	Decision string `form:"decision" json:"decision"`
}

func (req *Decision) ToDomain() *transactions.Domain {
	return &transactions.Domain{
		Status: req.Decision,
	}
}
