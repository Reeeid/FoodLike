package main

import (
	"log"
	"os"

	"foodlike-backend/internal/infrastructure/db"
	"foodlike-backend/internal/infrastructure/router"
)

func main() {
	conn, err := db.NewMySQLConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.New(conn)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
