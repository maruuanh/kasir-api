package handler

import (
	"encoding/json"
	"kasir-api/models"
	"net/http"
	"strconv"
	"strings"
)

// GET /api/categories/{id}
// @Summary Get Categories by ID
// @Description Mengambil data categories berdasarkan ID
// @Param id path int true "Categories ID"
// @Success 200 {object} models.Categories
// @Failure 400 {string} string "Invalid Categories ID"
// @Failure 404 {string} string "Categories tidak ditemukan"
// @Router /api/categories/{id} [get]
func GetCategoriesByID(w http.ResponseWriter, idStr string, categories []models.Categories) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Categories ID", http.StatusBadRequest)
		return
	}

	for _, p := range categories {
		if p.ID == id {
			JsonResponse(w, p)
			return
		}
	}

	http.Error(w, "Categories tidak ditemukan", http.StatusNotFound)
}

// PUT /api/categories/{id}
// @Summary Update Categories by ID
// @Description Memperbarui data categories berdasarkan ID
// @Param id path int true "Categories ID"
// @Param categories body models.Categories true "Updated Categories Data"
// @Success 200 {object} models.Categories
// @Failure 400 {string} string "Invalid request"
// @Router /api/categories/{id} [put]
func UpdateCategories(w http.ResponseWriter, r *http.Request, categories *[]models.Categories) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Categories ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updatedCategories models.Categories
	err = json.NewDecoder(r.Body).Decode(&updatedCategories)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop categories, cari id, ganti sesuai data dari request
	for i := range *categories {
		if (*categories)[i].ID == id {
			updatedCategories.ID = id
			(*categories)[i] = updatedCategories

			JsonResponse(w, updatedCategories)
			return
		}
	}
}

// DELETE /api/categories/{id}
// @Summary Delete Categories by ID
// @Description Menghapus data categories berdasarkan ID
// @Param id path int true "Categories ID"
// @Success 200 {string} string "Categories berhasil dihapus"
// @Failure 400 {string} string "Invalid Categories ID"
// @Router /api/categories/{id} [delete]
func DeleteCategories(w http.ResponseWriter, r *http.Request, categories *[]models.Categories) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	// ganti ke int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Categories ID", http.StatusBadRequest)
		return
	}
	// loop categories cari ID,dapet index yang mau dihapus
	for i, p := range *categories {
		if p.ID == id {
			// buat slice baru dengan data sebelum dan sesudah index
			*categories = append((*categories)[:i], (*categories)[i+1:]...)
			JsonResponse(w, map[string]string{
				"message": "Categories berhasil dihapus",
			})
			return
		}
	}
	// bikin slice baru dengan data sebelum dan sesudah index
}

// POST /api/categories/
// @Summary Create new Categories
// @Description Menambahkan data categories baru
// @Param categories body models.Categories true "New Categories Data"
// @Success 201 {object} models.Categories
// @Failure 400 {string} string "Invalid request"
// @Router /api/categories/ [post]
func PostCategories(w http.ResponseWriter, r *http.Request, categories *[]models.Categories) {
	var categoriesBaru models.Categories
	err := json.NewDecoder(r.Body).Decode(&categoriesBaru)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	categoriesBaru.ID = len(*categories) + 1
	*categories = append(*categories, categoriesBaru)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(categoriesBaru)
}

// GET /api/categories
// @Summary Get all Categories
// @Description Mengambil semua data categories
// @Success 200 {array} models.Categories
// @Router /api/categories [get]
func GetAllCategories(w http.ResponseWriter, categories []models.Categories) {
	JsonResponse(w, categories)
}

func JsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
