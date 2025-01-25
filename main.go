package main

import (
	"log"
	"net/http"
	"os"
	"url-shortener/db"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := db.MongoDBConfig{
		URI:      os.Getenv("URI"),
		Database: os.Getenv("MONGO_DB"),
	}
	_, database, err := db.InitializeMongoDB(config)
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewRepository(database)
	service := service.NewURLService(repo)
	handler := handler.NewHandler(service)
	http.HandleFunc("/shorten", handler.ShortenURL)
	http.HandleFunc("/", handler.RedirectURL)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
