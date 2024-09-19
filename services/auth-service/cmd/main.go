package main

import (
	"log"
	"net/http"
	"os"

	"gopher-market/services/auth-service/internal/handlers"
	"gopher-market/services/auth-service/internal/repository"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	authRepo := repository.NewAuthRepository(redisClient)
	authHandler := handlers.NewAuthHandler(authRepo)

	r := mux.NewRouter()
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	log.Println("Starting Auth Service on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
