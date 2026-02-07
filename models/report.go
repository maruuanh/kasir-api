package models

type Report struct {
	TotalRevenue   float64 `json:"total_revenue"`
	TotalTransaksi int     `json:"total_transaksi"`
	ProdukTerlaris struct {
		Nama       string `json:"nama"`
		QtyTerjual int    `json:"qty_terjual"`
	} `json:"produk_terlaris"`
}
