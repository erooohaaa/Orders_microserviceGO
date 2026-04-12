package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"Orders/internal/app"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	paymentGRPCAddr := os.Getenv("PAYMENT_GRPC_ADDR") // Добавь в .env: PAYMENT_GRPC_ADDR=localhost:50051

	application := app.New(db, paymentGRPCAddr)

	// 1. Запуск REST API (для пользователей) в фоне
	go func() {
		log.Println("Order REST API listening on :8080")
		if err := application.HTTPServer.ListenAndServe(); err != nil {
			log.Fatalf("REST server error: %v", err)
		}
	}()

	// 2. Запуск gRPC сервера (для стриминга обновлений)
	grpcPort := "50052"
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen gRPC: %v", err)
	}

	log.Printf("Order gRPC Streaming Service listening on %v", grpcPort)
	if err := application.GRPCServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
