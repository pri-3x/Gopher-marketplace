package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"gopher-market/services/auth-service/internal/handlers"
	"gopher-market/services/auth-service/internal/repository"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	log.Printf("REDIS_ADDR environment variable: %s", redisAddr)

	if redisAddr == "" {
		redisAddr = "redis:6379"
		log.Printf("REDIS_ADDR was empty, using default: %s", redisAddr)
	}

	log.Printf("Attempting to connect to Redis at: %s", redisAddr)

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")

	authRepo := repository.NewAuthRepository(redisClient)
	authHandler := handlers.NewAuthHandler(authRepo)

	r := mux.NewRouter()
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	log.Println("Starting Auth Service on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
