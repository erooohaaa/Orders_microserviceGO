package app

import (
	"database/sql"
	"net/http"
	"time"

	transportHTTP "Orders/internal/transport/http"
	"Orders/internal/repository"
	"Orders/internal/usecase"
)

type App struct {
	Server *http.Server
}

func New(db *sql.DB, paymentServiceURL string) *App {
	
	orderRepo := repository.NewPostgresOrderRepository(db)

	httpClient := &http.Client{Timeout: 2 * time.Second}
	paymentGateway := usecase.NewHTTPPaymentGateway(httpClient, paymentServiceURL)

	orderUC := usecase.NewOrderUseCase(orderRepo, paymentGateway)

	handler := transportHTTP.NewOrderHandler(orderUC)
	mux := http.NewServeMux()
	handler.Register(mux)

	return &App{
		Server: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}
}
