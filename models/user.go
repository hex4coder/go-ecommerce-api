package models

import "time"

type SoftDeleteTime struct {
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type SoftDeleteTimeWithoutDeletedAt struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type User struct {
	Id       int     `json:"id" gorm:"primary_key"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Role     int     `json:"role"`
	Password string  `json:"password"`
	Address  Address `json:"address,omitempty" gorm:"foreignKey:user_id;references:id"`
	SoftDeleteTimeWithoutDeletedAt
}

type Address struct {
	Id        int    `json:"id,omitempty" gorm:"primary_key"`
	UserID    int    `json:"user_id,omitempty"`
	Nomorhp   string `json:"nomorhp,omitempty"`
	Provinsi  string `json:"provinsi,omitempty"`
	Kota      string `json:"kota,omitempty"`
	Kecamatan string `json:"kecamatan,omitempty"`
	Desa      string `json:"desa,omitempty"`
	Dusun     string `json:"dusun,omitempty"`
	Jalan     string `json:"jalan,omitempty"`
	Kodepos   string `json:"kodepos,omitempty"`
	SoftDeleteTime
}
