package models

import "time"

type Order struct {
	Id               uint
	Tanggal          time.Time
	Status           string
	AlasanPembatalan string
	TotalHargaProduk uint64
	TotalDiskon      uint
	TotalBayar       uint64
	BuktiTransfer    string
	SudahTerbayar    bool
	CodePromo        string
	UserId           uint

	// soft deleted function
	SoftDeleteTime
}

type DetailOrder struct {
	Id uint
	OrderId uint
	ProductId uint
}
