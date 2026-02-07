package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetReport(start_date string, end_date string) ([]models.Report, error) {
	var report []models.Report
	var scanReport models.Report
	args := []interface{}{}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	dateFilter := " WHERE DATE(t.created_at) = CURRENT_DATE"
	if start_date != "" && end_date != "" {
		dateFilter = " WHERE DATE(t.created_at) BETWEEN $1 AND $2"
		args = append(args, start_date, end_date)
	}

	summaryQuery := "SELECT COALESCE(SUM(total_amount), 0), COUNT(id) FROM transactions t" + dateFilter
	err = r.db.QueryRow(summaryQuery, args...).Scan(&scanReport.TotalRevenue, &scanReport.TotalTransaksi)

	if err != nil {
		return nil, err
	}

	topProductQuery := `SELECT 
				p.name, 
				SUM(td.quantity) as qty_terjual 
			FROM transaction_details td
			JOIN products p ON td.product_id = p.id
			JOIN transactions t ON td.transaction_id = t.id ` + dateFilter + ` GROUP BY p.name
				ORDER BY qty_terjual DESC
				LIMIT 1`

	err = r.db.QueryRow(topProductQuery, args...).Scan(&scanReport.ProdukTerlaris.Nama, &scanReport.ProdukTerlaris.QtyTerjual)
	if err == sql.ErrNoRows {
		scanReport.ProdukTerlaris.Nama = "-"
		scanReport.ProdukTerlaris.QtyTerjual = 0
	} else if err != nil {
		return nil, err
	}

	report = append(report, scanReport)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return report, nil

}
