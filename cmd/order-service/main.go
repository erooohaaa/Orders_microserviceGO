package main

import (
	"database/sql"
	"log"
	"os"

	"Orders/internal/app"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}
	dsn := os.Getenv("DB_DSN")
	paymentURL := os.Getenv("PAYMENT_SERVICE_URL")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	log.Println("Order Service: connected to database")

	a := app.New(db, paymentURL)
	log.Println("Order Service: listening on :8080")
	if err := a.Server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
