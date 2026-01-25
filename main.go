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

func main() {
	// @title Kasir API
	// @version 1.0
	// @description API sederhana untuk manajemen categories di kasir

	var categories = models.DataCategories

	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
			if idStr == "" {
				handler.GetAllCategories(w, categories)
			} else {
				handler.GetCategoriesByID(w, idStr, categories)
			}
		case "POST":
			handler.PostCategories(w, r, &categories)
		case "PUT":
			handler.UpdateCategories(w, r, &categories)
		case "DELETE":
			handler.DeleteCategories(w, r, &categories)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
			return
		}

		http.NotFound(w, r)
	})

	fmt.Println("Server running di localhost:8080")

	// Jalankan server di port 8080
	err := http.ListenAndServe(":8080", nil)

	// Tangani error jika server gagal dijalankan
	if err != nil {
		fmt.Println("Gagal menjalankan server karena", err.Error())
	}
}
