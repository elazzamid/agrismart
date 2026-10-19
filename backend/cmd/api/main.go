package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elazzamid/agrismart/backend/internal/auth"
	"github.com/elazzamid/agrismart/backend/internal/farm"
	"github.com/elazzamid/agrismart/backend/internal/knowledge"
	"github.com/elazzamid/agrismart/backend/internal/platform"
)

const maxRequestBodyBytes int64 = 1 << 20 // 1 MiB

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
	farmHandler := farm.NewHandler(farmService)
	plotHandler := farm.NewPlotHandler(farm.NewPlotService(db))
	cropCycleHandler := farm.NewCropCycleHandler(farm.NewCropCycleService(db))
	catalogHandler := farm.NewCatalogHandler(farm.NewCatalogService(db))
	knowledgeHandler := knowledge.NewHTTPHandler(knowledge.NewService(db))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", healthHandler)
	mux.HandleFunc("GET /api/v1/ready", func(w http.ResponseWriter, r *http.Request) {
		readinessHandler(w, r, db)
	})
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.Handle("GET /api/v1/auth/me", authHandler.Authenticated(http.HandlerFunc(authHandler.Me)))

	protected := func(next http.Handler) http.Handler { return authHandler.Authenticated(next) }
	mux.Handle("GET /api/v1/farms", protected(http.HandlerFunc(farmHandler.List)))
	mux.Handle("POST /api/v1/farms", protected(http.HandlerFunc(farmHandler.Create)))
	mux.Handle("GET /api/v1/farms/{farmID}", protected(http.HandlerFunc(farmHandler.Get)))
	mux.Handle("GET /api/v1/farms/{farmID}/plots", protected(http.HandlerFunc(plotHandler.List)))
	mux.Handle("POST /api/v1/farms/{farmID}/plots", protected(http.HandlerFunc(plotHandler.Create)))
	mux.Handle("GET /api/v1/farms/{farmID}/plots/{plotID}/crop-cycles", protected(http.HandlerFunc(cropCycleHandler.List)))
	mux.Handle("POST /api/v1/farms/{farmID}/plots/{plotID}/crop-cycles", protected(http.HandlerFunc(cropCycleHandler.Create)))
	mux.Handle("GET /api/v1/crops", protected(http.HandlerFunc(catalogHandler.ListCrops)))
	mux.Handle("GET /api/v1/crops/{cropID}/varieties", protected(http.HandlerFunc(catalogHandler.ListVarieties)))
	mux.Handle("GET /api/v1/crops/{cropID}/growth-stages", protected(http.HandlerFunc(catalogHandler.ListGrowthStages)))

	mux.Handle("GET /api/v1/knowledge", protected(http.HandlerFunc(knowledgeHandler.SearchPublished)))
	mux.Handle("GET /api/v1/knowledge/fertilizers/recommendations", protected(http.HandlerFunc(knowledgeHandler.RecommendFertilizers)))
	mux.Handle("POST /api/v1/diagnosis", protected(http.HandlerFunc(knowledgeHandler.Diagnose)))
	mux.Handle("POST /api/v1/knowledge/documents", protected(http.HandlerFunc(knowledgeHandler.CreateDocument)))
	mux.Handle("POST /api/v1/knowledge/documents/{id}/versions", protected(http.HandlerFunc(knowledgeHandler.AddVersion)))
	mux.Handle("POST /api/v1/knowledge/documents/{id}/validate", protected(http.HandlerFunc(knowledgeHandler.Validate)))
	mux.Handle("POST /api/v1/knowledge/documents/{id}/publish", protected(http.HandlerFunc(knowledgeHandler.Publish)))

	addr := os.Getenv("API_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           limitRequestBody(mux, maxRequestBodyBytes),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("AgriSmart API listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func limitRequestBody(next http.Handler, limit int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > limit {
			writeStatusJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "request body too large"})
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, limit)
		next.ServeHTTP(w, r)
	})
}

type databasePinger interface {
	Ping(context.Context) error
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeStatusJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func readinessHandler(w http.ResponseWriter, r *http.Request, db databasePinger) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		writeStatusJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready"})
		return
	}
	writeStatusJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func writeStatusJSON(w http.ResponseWriter, status int, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
