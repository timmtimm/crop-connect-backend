package users

import (
	"errors"
	"marketplace-backend/util"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepository Repository
}

func NewUserUseCase(ur Repository) UseCase {
	return &UserUseCase{
		userRepository: ur,
	}
}

/*
Create
*/

func (uu *UserUseCase) Register(domain *Domain) (string, int, error) {
	isRoleAvailable := util.CheckStringOnArray([]string{"buyer", "farmer"}, domain.Role)
	if !isRoleAvailable {
		return "", http.StatusBadRequest, errors.New("role tidak valid")
	}

	_, err := uu.userRepository.GetByEmail(domain.Email)
	if err == nil {
		return "", http.StatusConflict, errors.New("email telah terdaftar")
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(domain.Password), bcrypt.DefaultCost)
	domain.ID = primitive.NewObjectID()
	domain.Password = string(encryptedPassword)
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	user, err := uu.userRepository.Create(domain)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	token := util.GenerateToken(user.ID.Hex(), user.Role)
	return token, http.StatusCreated, nil
}

/*
Read
*/

func (uu *UserUseCase) Login(domain *Domain) (string, int, error) {
	user, err := uu.userRepository.GetByEmail(domain.Email)
	if err == mongo.ErrNoDocuments {
		return "", http.StatusNotFound, errors.New("email tidak terdaftar")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(domain.Password))
	if err != nil {
		return "", http.StatusUnauthorized, errors.New("password salah")
	}

	token := util.GenerateToken(user.ID.Hex(), user.Role)
	return token, http.StatusOK, nil
}

func (uu *UserUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	user, err := uu.userRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("gagal mendapatkan user")
	}

	return user, nil
}

/*
Update
*/

func (uu *UserUseCase) UpdateProfile(domain *Domain) (Domain, error) {
	user, err := uu.userRepository.GetByID(domain.ID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, errors.New("user tidak ditemukan")
	}

	if domain.Email != user.Email {
		_, err := uu.userRepository.GetByEmail(domain.Email)
		if err == nil {
			return Domain{}, errors.New("email telah terdaftar")
		}
	}

	user.Name = domain.Name
	user.Description = domain.Description
	user.Email = domain.Email
	user.PhoneNumber = domain.PhoneNumber
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	user, err = uu.userRepository.Update(&user)
	if err != nil {
		return Domain{}, errors.New("gagal mengupdate user")
	}

	return user, nil
}

/*
Delete
*/
