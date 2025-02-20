package controllers

import (
	"fmt"
	"sync"
	"time"

	"github.com/hex4coder/go-ecommerce-api/models"
	"gorm.io/gorm"
)

type NewOrderRequest struct {
	// user id
	UserId uint `json:"user_id"`

	// promo code
	CodePromo string `json:"code_promo"`

	// order
	TotalHargaProduk uint64 `json:"total_harga_produk"`
	TotalDiskon      uint   `json:"total_diskon"`
	TotalBayar       uint64 `json:"total_bayar"`
	BuktiTransfer    string `json:"bukti_transfer"`

	// jumlah barang yang dibeli
	Detail []OrderDetailRequest `json:"detail"`
}

type CancelOrderRequest struct {
	AlasanPembatalan string `json:"alasan_pembatalan,omitempty"`
	Id               uint   `json:"id"`
}

type OrderDetailRequest struct {
	ProductId  uint   `json:"product_id"`
	Jumlah     uint   `json:"jumlah"`
	Harga      uint   `json:"harga"`
	Total      uint   `json:"total"`
	Ukuran     string `json:"ukuran"`
	Keterangan string `json:"keterangan,omitempty"`
}

// create order api
type OrderAPI struct {
	db *gorm.DB
}

func NewOrderAPI(db *gorm.DB) *OrderAPI {
	return &OrderAPI{
		db: db,
	}
}

// implements all of the interface
func (o *OrderAPI) GetMyOrders(userId int) ([]*models.Order, error) {

	// create the array/slice
	orders := []*models.Order{}

	// find the orders corresponding to user id
	r := o.db.Table("pesanan").Where("user_id = ?", userId).Find(orders)

	// check if there is an error
	if r.Error != nil {
		return nil, r.Error
	}

	// success, return the data
	return orders, nil
}

func (o *OrderAPI) GetDetailOrder(orderId int) (*models.Order, []*models.DetailOrder, error) {
	// create the detail order and order vars
	var (
		order   *models.Order
		details []*models.DetailOrder
	)

	// find order with id
	r := o.db.Table("pesanan").Where("id = ?", orderId).First(order)

	// check error
	if r.Error != nil {
		return nil, nil, r.Error
	}

	// check if the order is not found
	if r.RowsAffected < 1 {
		return nil, nil, fmt.Errorf("no order with id %d", orderId)
	}

	// order is found then, find the items order
	d := o.db.Table("detail_pesanan").Where("pesanan_id = ?", orderId).Find(details)

	// check error
	if d.Error != nil {
		return nil, nil, d.Error
	}

	// no error, return the data
	return order, details, nil
}

func (o *OrderAPI) GetOrderStatus(orderId int) (string, error) {

	// create the template variable
	order := new(models.Order)

	// find the data
	r := o.db.Table("pesanan").Where("id = ?", orderId).First(order)

	// check error
	if r.Error != nil {
		return "", r.Error
	}

	// check if the order is not found
	if r.RowsAffected < 1 {
		return "", fmt.Errorf("no order with id %d", orderId)
	}

	// success, return the status
	return order.Status, nil
}

func (o *OrderAPI) CreateOrder(request *NewOrderRequest) error {

	// create template order and detail items order
	order := new(models.Order)

	// new order
	order.Status = "baru"
	order.Tanggal = time.Now()
	order.AlasanPembatalan = ""
	order.SudahTerbayar = false // default

	// file bukti transfer
	// terima file dari client
	order.BuktiTransfer = request.BuktiTransfer

	// fill order with data
	order.UserId = request.UserId
	order.CodePromo = request.CodePromo

	// calculate the discount and the sum of products
	order.TotalHargaProduk = request.TotalHargaProduk
	order.TotalDiskon = request.TotalDiskon
	order.TotalBayar = request.TotalBayar

	//TODO: Create new order with new detail order
	r := o.db.Table("pesanan").Create(order)

	// check error
	if r.Error != nil {
		return r.Error
	}

	// get the inserted id from query
	orderId := order.Id

	// create async process with sync.WaitGroup
	var wg *sync.WaitGroup
	for _, item := range request.Detail {
		wg.Add(1)

		// convert item to models.DetailOrder
		detail := &models.DetailOrder{
			PesananId:  orderId,
			ProdukId:   item.ProductId,
			Jumlah:     item.Jumlah,
			Harga:      uint64(item.Harga),
			Total:      uint64(item.Total),
			Keterangan: item.Keterangan,
			Ukuran:     item.Ukuran,
		}

		// create go routine for inserting data to table
		go func(db *gorm.DB, item *models.DetailOrder, wg *sync.WaitGroup) {
			defer wg.Done()

			// insert to database
			d := db.Table("detail_pesanan").Create(item)
			if d.Error == nil {

				// buat vars
				produk := new(models.Product)

				// kurangi stok produk dengan yang baru
				r := db.Table("produk").Where("id = ?", item.ProdukId).First(produk)

				if r.RowsAffected > 0 {
					// update stock
					db.Table("produk").Where("id = ?", item.ProdukId).Update("stok", produk.Stok-int(item.Jumlah))
				}

			}

		}(o.db, detail, wg)
	}

	// wait the process
	wg.Wait()

	// no error
	return nil
}

func (o *OrderAPI) CancelOrder(request *CancelOrderRequest) error {

	// template order
	order := new(models.Order)

	// find the order
	r := o.db.Table("pesanan").Where("id = ?", request.Id).First(order)

	// check error
	if r.Error != nil {
		return r.Error
	}

	// check if order is found
	if r.RowsAffected < 1 {
		return fmt.Errorf("no order with id %d", request.Id)
	}

	// check if status is new
	if order.Status == "baru" || order.Status == "sedang diproses" {

		// semua item dikembalikan stoknya
		details := []*models.DetailOrder{}

		// find the details
		d := o.db.Table("detail_pesanan").Where("pesanan_id = ?", request.Id).Find(details)

		// if error
		if d.Error != nil {
			return d.Error
		}

		// update all the stock
		for _, item := range details {
			// kembalikan stock dari produk yang dibeli
			product := new(models.Product)

			p := o.db.Table("produk").Where("id = ?", item.ProdukId).First(product)

			// if no error and there is a product
			if p.Error == nil && p.RowsAffected > 0 {
				product.Stok = product.Stok + int(item.Jumlah)
				o.db.Table("produk").Where("id = ?", item.ProdukId).Updates(product)
			}
		}

		// update the order status and the reason
		order.Status = "dibatalkan"
		order.AlasanPembatalan = request.AlasanPembatalan

		// commit the changes to database
		c := o.db.Table("pesanan").Where("id = ?", request.Id).Updates(order)

		// check error
		if c.Error != nil {
			return c.Error
		}

		// no error
		return nil
	} else {
		// error beacause the status is completed or has been sended
		return fmt.Errorf("tidak bisa membatalkan pesanan karena statusnya : %s", order.Status)
	}
}

func (o *OrderAPI) DeleteOrder(orderId int) error {

	// vars
	var (
		order   *models.Order
		details []*models.DetailOrder
	)

	// find the order
	r := o.db.Table("pesanan").Where("id = ?", orderId).First(order)

	// check if error
	if r.Error != nil {
		return r.Error
	}

	// if no record
	if r.RowsAffected < 1 {
		return fmt.Errorf("pesanan dengan id %d tidak ditemukan", orderId)
	}

	// delete the record
	// 1. delete all the items or detail orders
	// find the details
	d := o.db.Table("detail_pesanan").Where("pesanan_id = ?", orderId).Find(&details)

	// check error
	if d.Error != nil {
		return d.Error
	}

	// if there are datas
	if d.RowsAffected > 0 {
		// create an asynchronous operations
		var wg *sync.WaitGroup

		// process
		for _, item := range details {
			// wait
			wg.Add(1)

			// delete using concurrency
			go func(wg *sync.WaitGroup, db *gorm.DB, detail *models.DetailOrder) {
				// make sure its called
				defer wg.Done()

				// delete the order items
				db.Table("detail_pesanan").Delete(&models.DetailOrder{Id: detail.Id})
			}(wg, o.db, item)
		}

		// wait until finish
		wg.Wait()
	}

	// 2. delete the orders
	o.db.Table("pesanan").Delete(&models.Order{Id: uint(orderId)})

	// success and no error
	return nil
}
