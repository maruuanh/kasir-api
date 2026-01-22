package handler

import (
	"encoding/json"
	"kasir-api/models"
	"net/http"
	"strconv"
	"strings"
)

// GET localhost:8080/api/produk/{id}
// @Summary Get Produk by ID
// @Description Mengambil data produk berdasarkan ID
// @Param id path int true "Produk ID"
// @Success 200 {object} models.Produk
// @Failure 400 {string} string "Invalid Produk ID"
// @Failure 404 {string} string "Produk tidak ditemukan"
// @Router /api/produk/{id} [get]
func GetProdukByID(w http.ResponseWriter, idStr string, produk []models.Produk) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			JsonResponse(w, p)
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// PUT localhost:8080/api/produk/{id}
// @Summary Update Produk by ID
// @Description Memperbarui data produk berdasarkan ID
// @Param id path int true "Produk ID"
// @Param produk body models.Produk true "Updated Produk Data"
// @Success 200 {object} models.Produk
// @Failure 400 {string} string "Invalid request"
// @Router /api/produk/{id} [put]
func UpdateProduk(w http.ResponseWriter, r *http.Request, produk *[]models.Produk) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updatedProduk models.Produk
	err = json.NewDecoder(r.Body).Decode(&updatedProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop produk, cari id, ganti sesuai data dari request
	for i := range *produk {
		if (*produk)[i].ID == id {
			updatedProduk.ID = id
			(*produk)[i] = updatedProduk

			JsonResponse(w, updatedProduk)
			return
		}
	}
}

// DELETE localhost:8080/api/produk/{id}
// @Summary Delete Produk by ID
// @Description Menghapus data produk berdasarkan ID
// @Param id path int true "Produk ID"
// @Success 200 {string} string "Produk berhasil dihapus"
// @Failure 400 {string} string "Invalid Produk ID"
// @Router /api/produk/{id} [delete]
func DeleteProduk(w http.ResponseWriter, r *http.Request, produk *[]models.Produk) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	// ganti ke int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}
	// loop produk cari ID,dapet index yang mau dihapus
	for i, p := range *produk {
		if p.ID == id {
			// buat slice baru dengan data sebelum dan sesudah index
			*produk = append((*produk)[:i], (*produk)[i+1:]...)
			JsonResponse(w, map[string]string{
				"message": "Produk berhasil dihapus",
			})
			return
		}
	}
	// bikin slice baru dengan data sebelum dan sesudah index
}

// POST localhost:8080/api/produk/
// @Summary Create new Produk
// @Description Menambahkan data produk baru
// @Param produk body models.Produk true "New Produk Data"
// @Success 201 {object} models.Produk
// @Failure 400 {string} string "Invalid request"
// @Router /api/produk/ [post]
func PostProduk(w http.ResponseWriter, r *http.Request, produk *[]models.Produk) {
	var produkBaru models.Produk
	err := json.NewDecoder(r.Body).Decode(&produkBaru)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	produkBaru.ID = len(*produk) + 1
	*produk = append(*produk, produkBaru)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(produkBaru)
}

// GET localhost:8080/api/produk
// @Summary Get all Produk
// @Description Mengambil semua data produk
// @Success 200 {array} models.Produk
// @Router /api/produk [get]
func GetAllProduk(w http.ResponseWriter, produk []models.Produk) {
	JsonResponse(w, produk)
}

func JsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
