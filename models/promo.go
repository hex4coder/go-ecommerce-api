package models

type PromoCode struct {
	Code     string `json:"code" gorm:"unique"`
	Type     string `json:"type"`
	Discount int    `json:"discount"`

	SoftDeleteTimeWithoutDeletedAt
}
