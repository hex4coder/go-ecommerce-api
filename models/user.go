package models


type User struct {
	Id       int       `json:"id" gorm:"primary_key"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     int       `json:"role"`
	Password string    `json:"password"`
	Address  Addresses `json:"address" gorm:"foreignKey:user_id;references:id"`
}

type Addresses struct {
	Id        int    `json:"id" gorm:"primary_key"`
	UserID    int    `json:"user_id"`
	NomorHP   string `json:"nomorhp"`
	Provinsi  string `json:"provinsi"`
	Kota      string `json:"kota"`
	Kecamatan string `json:"kecamatan"`
	Desa      string `json:"desa"`
	Dusun     string `json:"dusun"`
	Jalan     string `json:"jalan"`
	KodePos   string `json:"kodepos"`
}

