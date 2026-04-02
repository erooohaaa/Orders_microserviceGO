package usecase

import (
	"Orders/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrGatewayUnavailable = errors.New("payment service unavailable")

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, id string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type PaymentGateway interface {
	Authorize(ctx context.Context, orderID string, amount int64) (transactionID string, status string, err error)
}

type OrderUseCase struct {
	repo    OrderRepository
	payment PaymentGateway
}

func NewOrderUseCase(repo OrderRepository, payment PaymentGateway) *OrderUseCase {
	return &OrderUseCase{repo: repo, payment: payment}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, customerID, itemName string, amount int64) (*domain.Order, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	order := &domain.Order{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     domain.StatusPending,
	}

	if err := uc.repo.Save(ctx, order); err != nil {
		return nil, err
	}

	_, payStatus, err := uc.payment.Authorize(ctx, order.ID, order.Amount)
	if err != nil {

		_ = uc.repo.UpdateStatus(ctx, order.ID, domain.StatusFailed)
		order.Status = domain.StatusFailed
		return nil, ErrGatewayUnavailable
	}

	newStatus := domain.StatusFailed
	if payStatus == "Authorized" {
		newStatus = domain.StatusPaid
	}
	_ = uc.repo.UpdateStatus(ctx, order.ID, newStatus)
	order.Status = newStatus
	return order, nil
}

func (uc *OrderUseCase) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *OrderUseCase) CancelOrder(ctx context.Context, id string) error {
	order, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if !order.CanBeCancelled() {
		return errors.New("only Pending orders can be cancelled")
	}
	return uc.repo.UpdateStatus(ctx, id, domain.StatusCancelled)
}
