package controllers

import (
	"fmt"

	"github.com/hex4coder/go-ecommerce-api/models"
	"gorm.io/gorm"
)

type UserAPI struct {
	DB *gorm.DB
}

func (u *UserAPI) GetUsers() ([]models.User, error) {
	var users []models.User

	a := u.DB.Find(&users)
	if a.Error != nil {
		return nil, a.Error
	}
	return users, nil
}

func (u *UserAPI) GetUsersWithAddress() ([]models.User, error) {
	var users []models.User

	a := u.DB.Preload("Address").Find(&users)
	if a.Error != nil {
		return nil, a.Error
	}
	return users, nil
}

func (u *UserAPI) GetUserById(id int) (models.User, error) {

	var user models.User

	a := u.DB.Where(&models.User{Id: id}).Find(&user)
	if a.Error != nil {
		return models.User{}, a.Error
	}

	if a.RowsAffected < 1 {
		return models.User{}, fmt.Errorf("no user found with id %d", id)
	}

	return user, nil
}

func (u *UserAPI) GetUserAddressById(id int) (models.Addresses, error) {

	var address models.Addresses

	a := u.DB.Where(&models.Addresses{UserID: id}).Find(&address)
	if a.Error != nil {
		return models.Addresses{}, a.Error
	}

	if a.RowsAffected < 1 {
		return models.Addresses{}, fmt.Errorf("no address found with user id %d", id)
	}

	return address, nil
}
