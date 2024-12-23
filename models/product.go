package models

type Brand struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Logo string `json:"logo,omitempty"`

	SoftDeleteTime
}

type Kategori struct {
	Id           int    `json:"id"`
	NamaKategori string `json:"nama_kategori"`
	Slug         string `json:"slug"`
	Gambar       string `json:"gambar,omitempty"`
	SoftDeleteTime
}

type Product struct {
	Id        int     `json:"id"`
	Nama      string  `json:"nama"`
	Deskripsi string  `json:"deskripsi"`
	Harga     float64 `json:"harga"`
	Stok      int     `json:"stok"`
	IsPopular bool    `json:"is_popular"`
	Thumbnail string  `json:"thumbnail"`

	// foreign key
	KategoriID int `json:"kategori_id"`
	BrandID    int `json:"brand_id"`

	SoftDeleteTime
}

type PhotoProducts struct {
	Id       int     `json:"id"`
	ProdukID int     `json:"produk_id"`
	Foto     string  `json:"foto"`
	Produk   Product `json:"produk,omitempty"`
	SoftDeleteTime
}

type UkuranProduks struct {
	Id       int     `json:"id"`
	ProdukID int     `json:"produk_id"`
	Ukuran   string  `json:"ukuran"`
	Produk   Product `json:"produk,omitempty"`
	SoftDeleteTime
}
