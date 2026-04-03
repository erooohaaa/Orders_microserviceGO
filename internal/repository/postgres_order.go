package repository

import (
	"context"
	"database/sql"

	"Orders/internal/domain"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Save(ctx context.Context, o *domain.Order) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO orders (id, customer_id, item_name, amount, status)
         VALUES ($1, $2, $3, $4, $5)`,
		o.ID, o.CustomerID, o.ItemName, o.Amount, o.Status,
	)
	return err
}

func (r *PostgresOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, customer_id, item_name, amount, status, created_at
         FROM orders WHERE id = $1`, id)

	o := &domain.Order{}
	err := row.Scan(&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return o, err
}

func (r *PostgresOrderRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status = $1 WHERE id = $2`, status, id)
	return err
}
