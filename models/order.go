package models

import "time"

type Order struct {
	Id               uint
	UserId           uint
	Tanggal          time.Time
	TotalHargaBarang uint64
	TotalDiskon      uint
	TotalBayar       uint64
	Status           string
	SudahTerbayar    bool
	BuktiTransfer    string

	// soft deleted function
	SoftDeleteTime
}
