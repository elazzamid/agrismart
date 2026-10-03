package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/elazzamid/agrismart/backend/internal/auth"
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", healthHandler)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.Handle("GET /api/v1/auth/me", authHandler.Authenticated(http.HandlerFunc(authHandler.Me)))

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
