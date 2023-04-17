package request

import (
	"crop_connect/business/batchs"
	"crop_connect/helper"
	"errors"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

type Cancel struct {
	CancelReason string `form:"cancelReason" json:"cancelReason" validate:"required"`
}

func (req *Cancel) ToDomain() *batchs.Domain {
	return &batchs.Domain{
		CancelReason: req.CancelReason,
	}
}

func (req *Cancel) Validate() []helper.ValidationError {
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
