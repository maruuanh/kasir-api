package main

import (
	"fmt"
	"kasir-api/database"
	_ "kasir-api/docs"
	"kasir-api/handlers"
	"kasir-api/middlewares"
	"kasir-api/repositories"
	"kasir-api/services"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
	APIKey string `mapstructure:"API_KEY"`
}

func main() {
	// @title Kasir API
	// @version 1.0.1
	// @description API untuk aplikasi manajemen kasir yang di-update dengan menggunakan database PostgreSQL. Terdapat penambahan endpoint untuk mengelola kategori produk serta relasi antara produk dan kategori.

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
		APIKey: viper.GetString("API_KEY"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		fmt.Println("Gagal koneksi ke database:", err.Error())
		return
	}

	defer db.Close()

	apiKeyMiddleware := middlewares.APIkey(config.APIKey)

	// var categories = models.DataCategories
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	http.HandleFunc("/api/kategori", middlewares.CORS(middlewares.Logger(productHandler.GetCategories)))
	http.HandleFunc("/api/produk", middlewares.CORS(middlewares.Logger(productHandler.HandleProducts)))
	http.HandleFunc("/api/produk/", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(productHandler.HandleProductByID))))
	http.HandleFunc("/api/report/", middlewares.CORS(middlewares.Logger(reportHandler.HandleReport)))
	http.HandleFunc("/api/checkout", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(transactionHandler.HandleCheckout))))

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
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running at", addr)

	err = http.ListenAndServe(addr, nil)

	// Tangani error jika server gagal dijalankan
	if err != nil {
		fmt.Println("Gagal menjalankan server karena", err.Error())
	}
}
