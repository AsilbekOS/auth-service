package main

import (
	"auth-service/internal/auth"
	"auth-service/internal/db"
	"auth-service/pkg/config"
	"auth-service/pkg/email"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.NewConfig()
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	emailer := email.NewEmailer(cfg)
	authService := auth.NewService(database, cfg, emailer)

	database.Migrate()

	handler := auth.NewHandler(authService)
	r := mux.NewRouter()
	r.HandleFunc("/tokens/{userID}", handler.GenerateTokensHandler).Methods("POST")
	r.HandleFunc("/refresh", handler.RefreshTokenHandler).Methods("POST")

	log.Println("Service is listening to port:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
