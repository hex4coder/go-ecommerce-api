package controllers

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hex4coder/go-ecommerce-api/models"
	"github.com/joho/godotenv"
)

func TestGetKategori(t *testing.T) {
	// open database
	db, err := models.ConnectDB()

	if err != nil {
		t.Fatalf("%v", err)
	}

	start := time.Now()

	var wg sync.WaitGroup

	ec1 := make(chan error)
	dt1 := make(chan []*models.Brand)
	// defer close(ec1)
	// defer close(dt1)
	ec2 := make(chan error)
	dt2 := make(chan []*models.Kategori)
	// defer close(ec2)
	// defer close(dt2)
	wg.Add(2)

	go func(wg *sync.WaitGroup, ec chan error, dt chan []*models.Brand) {
		defer wg.Done()
		b := NewBrandAPI(db)
		data, err := b.GetAll()
		if err != nil {
			ec <- err
		}
		dt <- data
	}(&wg, ec1, dt1)

	go func(wg *sync.WaitGroup, ec chan error, dt chan []*models.Kategori) {
		defer wg.Done()
		b := NewKategoriAPI(db)
		data, err := b.GetAll()
		if err != nil {
			ec <- err
		}
		dt <- data
	}(&wg, ec2, dt2)

	for i := 0; i < 2; i++ {
		select {
		case e := <-ec1:
			fmt.Println(e.Error())
		case d := <-dt1:
			fmt.Println(d)
		case e := <-ec2:
			fmt.Println(e.Error())
		case d := <-dt2:
			fmt.Println(d)
		}
	}
	wg.Wait()
	end := time.Since(start)
	fmt.Println("Waktu kerja ", end)

}

func TestDetailProduct(t *testing.T) {

	// load env
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}

	// open database
	db, err := models.ConnectDB()

	if err != nil {
		t.Fatalf("%v", err)
	}

	// api
	b := NewBrandAPI(db)
	k := NewKategoriAPI(db)
	p := NewProductAPI(db)

	id := 1
	// buat temporary data
	data := make(map[string]any)

	// find products by categori id
	product, err := p.GetDetailProduct(id)
	if err != nil {
		t.Fatalf("Error get detail product %v\n", err)
		return
	}

	fmt.Println("Detail product", product)

	// assign product to data
	data["product"] = product

	// cari brand dari product tersebut
	brand, err := b.GetById(product.BrandID)
	if err != nil {
		t.Fatalf("Error get brand by id %v\n", err)
		return
	}

	// assign brand to data
	data["brand"] = brand

	// cari kategori detail dari product
	cat, err := k.GetById(product.KategoriID)
	if err != nil {
		t.Fatalf("Error get kategori by id %v\n", err)
		return
	}

	// assign category to data
	data["kategori"] = cat

	fmt.Println("list data", data)
}
