package controllers

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hex4coder/go-ecommerce-api/models"
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
