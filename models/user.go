package models

type User struct {
	Id       int       `json:"id" gorm:"primary_key"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     int       `json:"role"`
	Password string    `json:"password"`
	Address  Addresses `json:"address,omitempty" gorm:"foreignKey:user_id;references:id"`
}

type Addresses struct {
	Id        int    `json:"id,omitempty" gorm:"primary_key"`
	UserID    int    `json:"user_id,omitempty"`
	NomorHP   string `json:"nomorhp,omitempty"`
	Provinsi  string `json:"provinsi,omitempty"`
	Kota      string `json:"kota,omitempty"`
	Kecamatan string `json:"kecamatan,omitempty"`
	Desa      string `json:"desa,omitempty"`
	Dusun     string `json:"dusun,omitempty"`
	Jalan     string `json:"jalan,omitempty"`
	KodePos   string `json:"kodepos,omitempty"`
}
