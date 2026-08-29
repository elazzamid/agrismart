package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/elazzamid/agrismart/backend/internal/auth"
	"github.com/elazzamid/agrismart/backend/internal/knowledge"
	"github.com/elazzamid/agrismart/backend/internal/platform"
)

func main() {
	ctx := context.Background()
	db, err := platform.OpenDatabase(ctx)
	if err != nil { log.Fatal(err) }
	defer db.Close()

	tokens, err := auth.NewTokenService()
	if err != nil { log.Fatal(err) }
	authService := auth.NewService(db, tokens)
	authHandler := auth.NewHandler(authService)
	knowledgeHandler := knowledge.NewHTTPHandler(knowledge.NewService(db))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", healthHandler)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.Handle("GET /api/v1/auth/me", authHandler.Authenticated(http.HandlerFunc(authHandler.Me)))

	protectedKnowledge := func(next http.Handler) http.Handler { return authHandler.Authenticated(next) }
	mux.Handle("GET /api/v1/knowledge", protectedKnowledge(http.HandlerFunc(knowledgeHandler.SearchPublished)))
	mux.Handle("GET /api/v1/knowledge/fertilizers/recommendations", protectedKnowledge(http.HandlerFunc(knowledgeHandler.RecommendFertilizers)))
	mux.Handle("POST /api/v1/diagnosis", protectedKnowledge(http.HandlerFunc(knowledgeHandler.Diagnose)))
	mux.Handle("POST /api/v1/knowledge/documents", protectedKnowledge(http.HandlerFunc(knowledgeHandler.CreateDocument)))
	mux.Handle("POST /api/v1/knowledge/documents/{id}/versions", protectedKnowledge(http.HandlerFunc(knowledgeHandler.AddVersion)))
	mux.Handle("POST /api/v1/knowledge/documents/{id}/validate", protectedKnowledge(http.HandlerFunc(knowledgeHandler.Validate)))
	mux.Handle("POST /api/v1/knowledge/documents/{id}/publish", protectedKnowledge(http.HandlerFunc(knowledgeHandler.Publish)))

	addr := os.Getenv("API_ADDR")
	if addr == "" { addr = ":8080" }
	server := &http.Server{Addr: addr, Handler: mux}
	log.Printf("AgriSmart API listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatal(err) }
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
