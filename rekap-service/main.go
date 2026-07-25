package main 

import (
	"log"
	"os"

	"rekap-service/internal/rekap"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Mengambil URL Akademik Service dari Environment Variable (diatur di docker-compose)
	akademikURL := os.Getenv("AKADEMIK_BASE_URL")
	if akademikURL == "" {
		akademikURL = "http://localhost:8080" // Fallback untuk testing lokal
	}

	// Dependency Injection
	client := rekap.NewHTTPAkademikClient(akademikURL)
	service := rekap.NewService(client)
	handler := rekap.NewHttpAdapter(service)

	// Routing
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "rekap-service"})
	})

	api := r.Group("/api/v1/rekap")
	{
		api.GET("/top-ipk", handler.TopIPK)
		api.GET("/jurusan", handler.PerJurusan)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Rekap Service berjalan di port :%s dan terhubung ke Akademik di %s", port, akademikURL)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan Rekap Service: %v", err)
	}
}