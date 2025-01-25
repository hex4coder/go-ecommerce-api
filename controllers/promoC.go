package controllers

import (
	"fmt"

	"github.com/hex4coder/go-ecommerce-api/models"
	"gorm.io/gorm"
)

type PromoCodeAPI struct {
	db *gorm.DB
}

func NewPromoCodeAPI(db *gorm.DB) *PromoCodeAPI {
	return &PromoCodeAPI{db: db}
}

func (p *PromoCodeAPI) CheckPromo(code string) (*models.PromoCode, error) {

	// temp variable
	promo := new(models.PromoCode)

	// find the promo with code
	r := p.db.Table("promo_codes").Where("code = ?", code).First(promo)

	// check error
	if r.Error != nil {
		return nil, r.Error
	}

	// check if found
	if r.RowsAffected < 1 {
		return nil, fmt.Errorf("tidak ada promo dengan code %s", code)
	}

	// success return data without error
	return promo, nil
}
