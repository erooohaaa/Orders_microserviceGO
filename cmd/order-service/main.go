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
		log.Println("No .env file found, using environment")
	}

	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	paymentGRPCAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	application := app.New(db, paymentGRPCAddr)

	go func() {
		log.Println("Order REST API listening on :8080")
		if err := application.HTTPServer.ListenAndServe(); err != nil {
			log.Fatalf("REST server error: %v", err)
		}
	}()

	grpcPort := os.Getenv("ORDER_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	log.Printf("Order gRPC streaming service listening on :%s", grpcPort)
	if err := application.GRPCServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
