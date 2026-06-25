// Command api es el entrypoint del servidor HTTP.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ema/fixture/backend/internal/database"
	"github.com/ema/fixture/backend/internal/handler"
	sqliteRepo "github.com/ema/fixture/backend/internal/repository/sqlite"
	"github.com/ema/fixture/backend/internal/router"
	"github.com/ema/fixture/backend/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicializar base de datos
	db, err := database.New("fixture.db")
	if err != nil {
		log.Fatalf("Error al inicializar DB: %v", err)
	}
	defer db.Close()

	// Repositorios
	categoriaRepo := sqliteRepo.NewCategoriaRepository(db)
	autoRepo := sqliteRepo.NewAutoRepository(db)
	fixtureRepo := sqliteRepo.NewFixtureRepository(db)

	// Servicios
	categoriaSvc := service.NewCategoriaService(categoriaRepo)
	autoSvc := service.NewAutoService(autoRepo, categoriaRepo)
	fixtureSvc := service.NewFixtureService(fixtureRepo, categoriaRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler()
	categoriaHandler := handler.NewCategoriaHandler(categoriaSvc)
	autoHandler := handler.NewAutoHandler(autoSvc)
	fixtureHandler := handler.NewFixtureHandler(fixtureSvc)

	r := router.New(healthHandler, authHandler, categoriaHandler, autoHandler, fixtureHandler)

	// Servir archivos estáticos (uploads)
	os.MkdirAll("uploads", 0755)
	r.Handle("/uploads/*", http.StripPrefix("/uploads/",
		http.FileServer(http.Dir("uploads"))))

	addr := ":" + port
	log.Printf("🚀 Servidor iniciado en http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
