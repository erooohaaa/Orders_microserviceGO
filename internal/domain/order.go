package domain

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("order not found")

type OrderStatus = string

const (
	StatusPending   OrderStatus = "Pending"
	StatusPaid      OrderStatus = "Paid"
	StatusFailed    OrderStatus = "Failed"
	StatusCancelled OrderStatus = "Cancelled"
)

type Order struct {
	ID         string
	CustomerID string
	ItemName   string
	Amount     int64
	Status     OrderStatus
	CreatedAt  time.Time
}

func (o *Order) CanBeCancelled() bool {
	return o.Status == StatusPending
}
