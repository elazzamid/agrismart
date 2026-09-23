package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/elazzamid/agrismart/backend/internal/auth"
	"github.com/elazzamid/agrismart/backend/internal/farm"
	"github.com/elazzamid/agrismart/backend/internal/platform"
)

func main() {
	ctx := context.Background()
	db, err := platform.OpenDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tokens, err := auth.NewTokenService()
	if err != nil {
		log.Fatal(err)
	}
	authService := auth.NewService(db, tokens)
	authHandler := auth.NewHandler(authService)
	farmService := farm.NewService(db)
	plotService := farm.NewPlotService(db)
	catalogService := farm.NewCatalogService(db)
	farmHandler := farm.NewHandler(farmService, plotService, catalogService)
	authenticated := authHandler.Authenticated

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", healthHandler)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.Handle("GET /api/v1/auth/me", authenticated(http.HandlerFunc(authHandler.Me)))
	mux.Handle("GET /api/v1/farms", authenticated(http.HandlerFunc(farmHandler.List)))
	mux.Handle("POST /api/v1/farms", authenticated(http.HandlerFunc(farmHandler.Create)))
	mux.Handle("GET /api/v1/farms/{id}", authenticated(http.HandlerFunc(farmHandler.Get)))
	mux.Handle("GET /api/v1/farms/{id}/plots", authenticated(http.HandlerFunc(farmHandler.ListPlots)))
	mux.Handle("POST /api/v1/farms/{id}/plots", authenticated(http.HandlerFunc(farmHandler.CreatePlot)))
	mux.Handle("GET /api/v1/crops", authenticated(http.HandlerFunc(farmHandler.ListCrops)))
	mux.Handle("GET /api/v1/crops/{id}", authenticated(http.HandlerFunc(farmHandler.GetCrop)))
	mux.Handle("GET /api/v1/crops/{id}/varieties", authenticated(http.HandlerFunc(farmHandler.ListVarieties)))
	mux.Handle("GET /api/v1/crops/{id}/growth-stages", authenticated(http.HandlerFunc(farmHandler.ListGrowthStages)))

	addr := os.Getenv("API_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("AgriSmart API listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
