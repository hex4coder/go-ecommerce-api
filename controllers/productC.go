package controllers

import (
	"fmt"

	"github.com/hex4coder/go-ecommerce-api/models"
	"gorm.io/gorm"
)

// ------------------------------KATEGORI API----------------------------------------
type KategoriAPI struct {
	DB *gorm.DB
}

func NewKategoriAPI(db *gorm.DB) *KategoriAPI {
	return &KategoriAPI{DB: db}
}

func (k *KategoriAPI) GetAll() ([]*models.Kategori, error) {

	data := []*models.Kategori{}

	q := k.DB.Table("kategori").Find(&data)

	if q.Error != nil {
		return nil, q.Error
	}

	return data, nil
}

func (k *KategoriAPI) GetById(id int) (*models.Kategori, error) {

	data := new(models.Kategori)

	q := k.DB.Table("kategori").Where(&models.Kategori{Id: id}).First(data)

	if q.Error != nil {
		return nil, q.Error
	}

	return data, nil
}

// ------------------------------BRAND API----------------------------------------
type BrandAPI struct {
	DB *gorm.DB
}

func NewBrandAPI(db *gorm.DB) *BrandAPI {
	return &BrandAPI{DB: db}
}

func (b *BrandAPI) GetAll() ([]*models.Brand, error) {
	data := []*models.Brand{}

	q := b.DB.Table("brands").Find(&data)

	if q.Error != nil {
		return nil, q.Error
	}

	return data, nil
}

func (b *BrandAPI) GetById(id int) (*models.Brand, error) {
	data := new(models.Brand)

	q := b.DB.Table("brands").Where(&models.Brand{
		Id: id,
	}).First(data)

	if q.Error != nil {
		return nil, q.Error
	}

	return data, nil
}

// ------------------------------PRODUCT API----------------------------------------

type ProductAPI struct {
	DB *gorm.DB
}

func NewProductAPI(db *gorm.DB) *ProductAPI {
	return &ProductAPI{
		DB: db,
	}
}

// ------------------------------INTERFACE METHODS------------------------------------
func (p *ProductAPI) GetAllProducts() ([]*models.Product, error) {

	data := []*models.Product{}

	// urutkan berdasarkan waktu
	r := p.DB.Table("produk").Order("created_at desc").Find(&data)

	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}
func (p *ProductAPI) GetProductsByCategoryID(katID int) ([]*models.Product, error) {
	data := []*models.Product{}

	// urutkan berdasarkan waktu
	r := p.DB.Table("produk").Where("kategori_id = ?", katID).Order("created_at desc").Find(&data)

	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}
func (p *ProductAPI) GetProductsByBrandID(brandID int) ([]*models.Product, error) {
	data := []*models.Product{}

	// urutkan berdasarkan waktu
	r := p.DB.Table("produk").Where("brand_id = ?", brandID).Order("created_at desc").Find(&data)

	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}
func (p *ProductAPI) GetDetailProduct(productID int) (*models.Product, error) {
	data := new(models.Product)

	// ambil detail product
	r := p.DB.Table("produk").Where("id = ?", productID).Order("created_at desc").First(data)

	if r.Error != nil {
		return nil, r.Error
	}

	if r.RowsAffected < 1 {
		return nil, fmt.Errorf("no product count with id %d", productID)
	}

	return data, nil
}

func (p *ProductAPI) GetProductPhotosByID(productID int) ([]*models.PhotoProducts, error) {
	data := []*models.PhotoProducts{}

	// urutkan berdasarkan waktu
	r := p.DB.Table("foto_produks").Where("produk_id = ?", productID).Order("created_at desc").Find(&data)

	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}
func (p *ProductAPI) GetUkuranProdukByID(productID int) ([]*models.UkuranProduks, error) {

	data := []*models.UkuranProduks{}

	// urutkan berdasarkan waktu
	r := p.DB.Table("ukuran_produks").Where("produk_id = ?", productID).Order("created_at desc").Find(&data)

	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}
func (p *ProductAPI) GetPopularProducts(limit int) ([]*models.Product, error) {
	data := []*models.Product{}

	var r *gorm.DB

	if limit > 0 {
		// urutkan berdasarkan waktu
		r = p.DB.Table("produk").Where("is_popular = ?", true).Order("created_at desc").Limit(limit).Find(&data)
	} else {
		r = p.DB.Table("produk").Where("is_popular = ?", true).Order("created_at desc").Find(&data)
	}

	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}
