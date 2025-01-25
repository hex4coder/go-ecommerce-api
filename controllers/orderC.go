package controllers

import "time"

type OrderRequest struct {
	// user id
	UserId uint `json:"user_id"`

	// jumlah barang yang dibeli
	Detail []OrderDetailRequest `json:"detail"`

	Tanggal          time.Time `json:"tanggal"`
	Status           string    `json:"status"`
	AlasanPembatalan string    `json:"alasan_pembatalan"`
	CodePromo        string    `json:"code_promo"`
}

type OrderDetailRequest struct {
	ProductId uint `json:"product_id"`
	Jumlah    uint `json:"jumlah"`
	Harga     uint `json:"harga"`
	Total     uint `json:"total"`
}
