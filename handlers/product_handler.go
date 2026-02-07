package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GET /api/kategori
// @Summary      Get All Categories
// @Description  Mengambil semua data kategori produk (Challange (Optional))
// @Tags         category
// @Accept       json
// @Tags         produk
// @Produce      json
// @Success      200  {array}   models.Categories
// @Failure      500  {string}  string "Failed to get categories"
// @Router       /api/kategori [get]
func (h *ProductHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories()
	if err != nil {
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GET /api/produk
// @Summary      Get All Products
// @Description  Mengambil semua data produk. Terdapat opsi untuk mendapatkan detail kategori produk
// @Accept       json
// @Tags         produk
// @Produce      json
// @Param        details  query     bool  false  "Tampilkan Detail Kategori Produk"
// @Param        name  	query     string false  "Tampilkan Detail Kategori Produk Berdasarkan Pencarian Nama"
// @Success      200      {array}   models.Product
// @Failure      500      {string}  string "Failed to get products"
// @Router       /api/produk [get]
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()
		fullStr := query.Get("details")
		if fullStr == "true" {
			h.GetAll(w, r, true)
		} else {
			h.GetAll(w, r, false)
		}
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /api/produk/{id}
// GetByID
// @Summary      Get Product by ID
// @Description  Mengambil data produk berdasarkan ID. Terdapat opsi untuk mendapatkan detail kategori produk
// @Accept       json
// @Tags         produk
// @Produce      json
// @Param        id       path      int   true   "Product ID"
// @Param        details  query     bool  false  "Tampilkan Detail Kategori Produk" default(false)
// @Success      200      {object}  models.Product
// @Failure      400      {string}  string "Invalid product ID"
// @Failure      404      {string}  string "Product not found"
// @Router       /api/produk/{id} [get]
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()
		fullStr := query.Get("details")
		if fullStr == "true" {
			h.GetByID(w, r, true)
		} else {
			h.GetByID(w, r, false)
		}
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request, details bool) {
	var products []models.Product
	var err error

	name := r.URL.Query().Get("name")
	if details {
		products, err = h.service.GetAllDetails(name)
	} else {
		products, err = h.service.GetAll(name)
	}

	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// POST /api/produk
// @Summary Create New Product
// @Description Menambahkan data produk baru, data yang perlu diisi: { category_id, name, price, stock }
// @Accept json
// @Tags   produk
// @Produce json
// @Param product body models.Product true "New Product Data"
// @Success 201 {object} models.Product
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create product"
// @Router /api/produk [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	err = h.service.Create(&product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request, details bool) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product *models.Product

	if details {
		product, err = h.service.GetDetailsByID(id)
	} else {
		product, err = h.service.GetByID(id)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// GET /api/produk/{id}?details=true
func (h *ProductHandler) GetDetailsByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	product, err := h.service.GetDetailsByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// PUT /api/produk/{id}
// @Summary Update Product by ID
// @Description Memperbarui data produk berdasarkan ID, data yang dapat diubah: { category_id, name, price, stock }
// @Accept json
// @Tags   produk
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Updated Product Data"
// @Success 200 {object} models.Product
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to update product"
// @Router /api/produk/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product.ID = id
	err = h.service.Update(&product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// DELETE /api/produk/{id}
// @Summary Delete Product by ID
// @Description Menghapus data produk berdasarkan ID
// @Param id path int true "Product ID"
// @Tags   produk
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid product ID"
// @Failure 500 {string} string "Failed to delete product"
// @Router /api/produk/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
	}

	err = h.service.Delete(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted successfully",
	})
}
