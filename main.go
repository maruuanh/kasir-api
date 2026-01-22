package main

import (
	"fmt"
	_ "kasir-api/docs"
	"kasir-api/handler"
	"kasir-api/models"
	"net/http"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// @title Kasir API
	// @version 1.0
	// @description API sederhana untuk manajemen produk di kasir
	// @host localhost:8080
	// @schemes http https

	var produk = models.DataProduk

	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
			if idStr == "" {
				handler.GetAllProduk(w, produk)
			} else {
				handler.GetProdukByID(w, idStr, produk)
			}
		case "POST":
			handler.PostProduk(w, r, &produk)
		case "PUT":
			handler.UpdateProduk(w, r, &produk)
		case "DELETE":
			handler.DeleteProduk(w, r, &produk)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	http.Handle("/", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))

	fmt.Println("Server running di localhost:8080")

	// Jalankan server di port 8080
	err := http.ListenAndServe(":8080", enableCORS(http.NewServeMux()))

	// Tangani error jika server gagal dijalankan
	if err != nil {
		fmt.Println("Gagal menjalankan server karena", err.Error())
	}
}
