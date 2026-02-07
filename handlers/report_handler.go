package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strings"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GET /api/report
// @Summary      Get Transaction Report By Selected Date
// @Description  Mengambil laporan data transaksi penjualan barang berdasarkan tanggal yang dipilih
// @Accept       json
// @Tags         report
// @Produce      json
// @Param        start_date  query     string  false  "Tanggal awal (Format: YYYY-MM-DD)" example(2026-01-01)
// @Param        end_date    query     string  false  "Tanggal akhir (Format: YYYY-MM-DD)" example(2026-02-01)
// @Success      200      {array}   models.Report
// @Failure      500      {string}  string "Failed to get report"
// @Router       /api/report [get]
func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetReport(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
}

// GET /api/report/hari-ini
// @Summary      Get Today's Transaction Report
// @Description  Mengambil laporan data transaksi penjualan barang khusus hari ini
// @Accept       json
// @Tags         report
// @Produce      json
// @Success      200      {array}   models.Report
// @Failure      500      {string}  string "Failed to get report"
// @Router       /api/report/hari-ini [get]
func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	var report []models.Report
	var err error
	var startDate, endDate string

	if strings.Contains(r.URL.Path, "/hari-ini") {
		startDate = ""
		endDate = ""
	} else {
		startDate = r.URL.Query().Get("start_date")
		endDate = r.URL.Query().Get("end_date")
	}

	report, err = h.service.GetReport(startDate, endDate)

	if err != nil {
		http.Error(w, "Failed to get report: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
