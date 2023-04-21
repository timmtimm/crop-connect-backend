package forgot_password

import (
	"crop_connect/business/users"
	"crop_connect/constant"
	"crop_connect/helper/mailgun"
	"crop_connect/util"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type ForgotPasswordUseCase struct {
	forgotPasswordRepository Repository
	userRepository           users.Repository
	mailgun                  mailgun.Function
}

func NewUseCase(fpr Repository, ur users.Repository, mg mailgun.Function) UseCase {
	return &ForgotPasswordUseCase{
		forgotPasswordRepository: fpr,
		userRepository:           ur,
		mailgun:                  mg,
	}
}

var errorResponse = errors.New("token tidak dapat digunakan")

/*
Create
*/

func (fpu *ForgotPasswordUseCase) Generate(email string) (int, error) {
	_, err := fpu.userRepository.GetByEmail(email)
	if err != nil {
		return http.StatusInternalServerError, errors.New("email tidak terdaftar")
	}

	domain := Domain{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Token:     util.GenerateUUID(),
		IsUsed:    false,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		ExpiredAt: primitive.NewDateTimeFromTime(time.Now().Add(24 * time.Hour)),
	}
	_, err = fpu.forgotPasswordRepository.Create(&domain)
	if err != nil {
		return http.StatusInternalServerError, errors.New("gagal membuat token")
	}

	_, _, err = fpu.mailgun.SendOneMailUsingTemplate("Lupa password Crop Connect?", constant.MailgunForgotPasswordTemplate, domain.Email, "", map[string]string{
		"token": domain.Token,
	})
	if err != nil {
		if err := fpu.forgotPasswordRepository.HardDelete(domain.ID); err != nil {
			return http.StatusInternalServerError, errors.New("gagal menghapus token")
		}

		return http.StatusInternalServerError, errors.New("gagal mengirim email")
	}

	return http.StatusCreated, nil
}

/*
Read
*/

func (fpu *ForgotPasswordUseCase) ValidateToken(token string) (int, error) {
	forgotPassword, err := fpu.forgotPasswordRepository.GetByToken(token)
	if err != nil {
		return http.StatusForbidden, errorResponse
	}

	if forgotPassword.IsUsed {
		return http.StatusForbidden, errorResponse
	} else if forgotPassword.ExpiredAt.Time().Before(time.Now()) {
		return http.StatusForbidden, errorResponse
	}

	return http.StatusOK, nil
}

/*
Update
*/

func (fpu *ForgotPasswordUseCase) ResetPassword(token string, password string) (int, error) {
	forgotPassword, err := fpu.forgotPasswordRepository.GetByToken(token)
	if err != nil {
		return http.StatusForbidden, errorResponse
	}

	user, err := fpu.userRepository.GetByEmail(forgotPassword.Email)
	if err != nil {
		return http.StatusForbidden, errorResponse
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusForbidden, errorResponse
	}

	forgotPassword.IsUsed = true
	forgotPassword.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err = fpu.forgotPasswordRepository.Update(&forgotPassword)
	if err != nil {
		return http.StatusForbidden, errorResponse
	}

	user.Password = string(encryptedPassword)
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err = fpu.userRepository.Update(&user)
	if err != nil {
		return http.StatusForbidden, errorResponse
	}

	return http.StatusOK, nil
}

/*
Delete
*/
