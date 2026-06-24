// Command api es el entrypoint del servidor HTTP.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ema/fixture/backend/internal/handler"
	"github.com/ema/fixture/backend/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	healthHandler := handler.NewHealthHandler()
	r := router.New(healthHandler)

	addr := ":" + port
	log.Printf("🚀 Servidor iniciado en http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
