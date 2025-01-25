package models

import "time"

type Order struct {
	Id               uint      `json:"id" gorm:"primaryKey"`
	Tanggal          time.Time `json:"tanggal"`
	Status           string    `json:"status"`
	AlasanPembatalan string    `json:"alasan_pembatalan,omitempty"`
	TotalHargaProduk uint64    `json:"total_harga_produk"`
	TotalDiskon      uint      `json:"total_diskon"`
	TotalBayar       uint64    `json:"total_bayar"`
	BuktiTransfer    string    `json:"bukti_transfer"`
	SudahTerbayar    bool      `json:"sudah_terbayar"`
	CodePromo        string    `json:"code_promo"`
	UserId           uint      `json:"user_id"`

	// soft deleted function
	SoftDeleteTime
}

type DetailOrder struct {
	Id         uint   `json:"id"`
	PesananId  uint   `json:"pesanan_id"`
	ProdukId   uint   `json:"produk_id"`
	Ukuran     string `json:"ukuran"`
	Keterangan string `json:"keterangan"`
	Harga      uint64 `json:"harga"`
	Jumlah     uint   `json:"jumlah"`
	Total      uint64 `json:"total"`

	// soft deleted function
	SoftDeleteTime
}
