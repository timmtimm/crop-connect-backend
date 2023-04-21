package users

import (
	"crop_connect/business/regions"
	"crop_connect/constant"
	"crop_connect/helper"
	"crop_connect/util"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepository   Repository
	regionRepository regions.Repository
}

func NewUseCase(ur Repository, rr regions.Repository) UseCase {
	return &UserUseCase{
		userRepository:   ur,
		regionRepository: rr,
	}
}

/*
Create
*/

func (uu *UserUseCase) Register(domain *Domain) (string, int, error) {
	isRoleAvailable := util.CheckStringOnArray([]string{constant.RoleBuyer, constant.RoleFarmer}, domain.Role)
	if !isRoleAvailable {
		return "", http.StatusBadRequest, errors.New("role tersedia hanya buyer dan farmer")
	}

	_, err := uu.regionRepository.GetByID(domain.RegionID)
	if err == mongo.ErrNoDocuments {
		return "", http.StatusNotFound, errors.New("daerah tidak ditemukan")
	} else if err != nil {
		return "", http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	_, err = uu.userRepository.GetByEmail(domain.Email)
	if err == mongo.ErrNoDocuments {
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(domain.Password), bcrypt.DefaultCost)
		domain.ID = primitive.NewObjectID()
		domain.Password = string(encryptedPassword)
		domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

		user, err := uu.userRepository.Create(domain)
		if err != nil {
			return "", http.StatusInternalServerError, errors.New("gagal melakuakn registrasi user")
		}

		token := helper.GenerateToken(user.ID.Hex(), user.Role)
		return token, http.StatusCreated, nil
	} else {
		return "", http.StatusConflict, errors.New("email telah terdaftar")
	}
}

func (uu *UserUseCase) RegisterValidator(domain *Domain) (string, int, error) {
	_, err := uu.regionRepository.GetByID(domain.RegionID)
	if err == mongo.ErrNoDocuments {
		return "", http.StatusNotFound, errors.New("daerah tidak ditemukan")
	} else if err != nil {
		return "", http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	_, err = uu.userRepository.GetByEmail(domain.Email)
	if err == mongo.ErrNoDocuments {
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(domain.Password), bcrypt.DefaultCost)
		domain.ID = primitive.NewObjectID()
		domain.Password = string(encryptedPassword)
		domain.Role = constant.RoleValidator
		domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

		user, err := uu.userRepository.Create(domain)
		if err != nil {
			return "", http.StatusInternalServerError, err
		}

		token := helper.GenerateToken(user.ID.Hex(), user.Role)
		return token, http.StatusCreated, nil
	} else {
		return "", http.StatusConflict, errors.New("email telah terdaftar")
	}
}

/*
Read
*/

func (uu *UserUseCase) Login(domain *Domain) (string, int, error) {
	user, err := uu.userRepository.GetByEmail(domain.Email)
	if err == mongo.ErrNoDocuments {
		return "", http.StatusNotFound, errors.New("email tidak terdaftar")
	} else if err != nil {
		return "", http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(domain.Password))
	if err != nil {
		return "", http.StatusUnauthorized, errors.New("password salah")
	}

	token := helper.GenerateToken(user.ID.Hex(), user.Role)
	return token, http.StatusOK, nil
}

func (uu *UserUseCase) GetByID(id primitive.ObjectID) (Domain, int, error) {
	user, err := uu.userRepository.GetByID(id)
	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("gagal mendapatkan user")
	}

	return user, http.StatusOK, nil
}

func (uu *UserUseCase) GetByPaginationAndQuery(query Query) ([]Domain, int, int, error) {
	users, totalData, err := uu.userRepository.GetByQuery(query)
	if err != nil {
		return []Domain{}, 0, http.StatusInternalServerError, errors.New("gagal mendapatkan user")
	}

	return users, totalData, http.StatusOK, nil
}

/*
Update
*/

func (uu *UserUseCase) UpdateProfile(domain *Domain) (Domain, int, error) {
	user, err := uu.userRepository.GetByID(domain.ID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("user tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mengambil data pengguna")
	}

	if domain.RegionID != user.RegionID {
		_, err = uu.regionRepository.GetByID(domain.RegionID)
		if err == mongo.ErrNoDocuments {
			return Domain{}, http.StatusNotFound, errors.New("daerah tidak ditemukan")
		} else if err != nil {
			return Domain{}, http.StatusInternalServerError, errors.New("gagal mengambil data proposal")
		}
	}

	if domain.Email != user.Email {
		_, err := uu.userRepository.GetByEmail(domain.Email)
		if err == nil {
			return Domain{}, http.StatusConflict, errors.New("email telah terdaftar")
		}
	}

	user.Name = domain.Name
	user.Description = domain.Description
	user.Email = domain.Email
	user.PhoneNumber = domain.PhoneNumber
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	user, err = uu.userRepository.Update(&user)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mengupdate user")
	}

	return user, http.StatusOK, nil
}

func (uu *UserUseCase) UpdatePassword(domain *Domain, newPassword string) (Domain, int, error) {
	user, err := uu.userRepository.GetByID(domain.ID)
	if err == mongo.ErrNoDocuments {
		return Domain{}, http.StatusNotFound, errors.New("user tidak ditemukan")
	} else if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mengambil data pengguna")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(domain.Password))
	if err != nil {
		return Domain{}, http.StatusUnauthorized, errors.New("password salah")
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = string(encryptedPassword)
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	user, err = uu.userRepository.Update(&user)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, errors.New("gagal mengupdate user")
	}

	return user, http.StatusOK, nil
}

/*
Delete
*/
