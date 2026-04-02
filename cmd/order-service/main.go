package main

import (
	"database/sql"
	"log"
	"os"

	"Orders/internal/app"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost user=postgres password=0000 dbname=order_db sslmode=disable"
	paymentURL := os.Getenv("PAYMENT_SERVICE_URL")
	if paymentURL == "" {
		paymentURL = "http://localhost:8081"
	}

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
