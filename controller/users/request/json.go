package request

import (
	"errors"
	"marketplace-backend/business/users"
	"marketplace-backend/helper"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

type Register struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=13,number"`
	Password    string `json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
	Role        string `json:"role" validate:"required"`
}

func (req *Register) ToDomain() *users.Domain {
	return &users.Domain{
		Name:        req.Name,
		Description: req.Description,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
		Role:        req.Role,
	}
}

func (req *Register) Validate() []helper.ValidationError {
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

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=!@#$%^&*,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
}

func (req *Login) ToDomain() *users.Domain {
	return &users.Domain{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (req *Login) Validate() []helper.ValidationError {
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
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=13,number"`
}

func (req *Update) ToDomain() *users.Domain {
	return &users.Domain{
		Name:        req.Name,
		Description: req.Description,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
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
