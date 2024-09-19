package main

import (
	"log"
	"net/http"
	"os"

	"gopher-market/services/user-service/internal/handlers"
	"gopher-market/services/user-service/internal/models"
	"gopher-market/services/user-service/internal/repository"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Connect to PostgreSQL
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto Migrate the schema
	db.AutoMigrate(&models.User{})

	repo := repository.NewUserRepository(db)
	handler := handlers.NewUserHandler(repo)

	r := mux.NewRouter()
	r.HandleFunc("/users", handler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users", handler.ListUsers).Methods("GET")

	log.Println("User Service starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
