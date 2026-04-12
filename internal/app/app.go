package app

import (
	"database/sql"
	"log"
	"net/http"

	"Orders/internal/repository"
	transportGRPC "Orders/internal/transport/grpc"
	transportHTTP "Orders/internal/transport/http"
	"Orders/internal/usecase"
	"github.com/erooohaaa/orders-generated/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	HTTPServer *http.Server
	GRPCServer *grpc.Server
}

func New(db *sql.DB, paymentGRPCAddr string) *App {
	// Настройка клиента для связи с Payment Service
	conn, err := grpc.Dial(paymentGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to payment service: %v", err)
	}

	orderRepo := repository.NewPostgresOrderRepository(db)

	// Используем новый gRPC гейтвей вместо старого HTTP
	paymentGateway := usecase.NewGRPCPaymentGateway(conn)
	orderUC := usecase.NewOrderUseCase(orderRepo, paymentGateway)

	// Настройка HTTP (REST)
	httpHandler := transportHTTP.NewOrderHandler(orderUC)
	mux := http.NewServeMux()
	httpHandler.Register(mux)

	// Настройка gRPC (Streaming)
	grpcServer := grpc.NewServer()
	api.RegisterOrderServiceServer(grpcServer, transportGRPC.NewOrderStreamHandler(orderUC))

	return &App{
		HTTPServer: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		GRPCServer: grpcServer,
	}
}
