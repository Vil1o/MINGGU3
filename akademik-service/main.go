package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"siakad-service/internal/config"
	"siakad-service/internal/siakad"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi db: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}
	log.Println("Koneksi ke PostgreSQL berhasil!")

	// 1. Inisialisasi Secondary Adapters (Outbound ke DB)
	mhsPgAdapter := siakad.NewMahasiswaRepository(db)
	nilaiPgAdapter := siakad.NewNilaiRepository(db)

	// 2. Inisialisasi Core Domain dengan menyuntikkan Secondary Adapters (Port Terhubung)
	coreService := siakad.NewService(mhsPgAdapter, nilaiPgAdapter)

	// 3. Inisialisasi Primary Adapter (Inbound dari HTTP) dengan menyuntikkan Core Domain
	httpAdapter := siakad.NewHttpAdapter(coreService)

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	// Pasang Routing
	httpAdapter.RegisterRoutes(&r.RouterGroup)

	log.Printf("Server berjalan di port :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}