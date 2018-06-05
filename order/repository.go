package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrderByID(ctx context.Context, id string) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert order
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	)
	if err != nil {
		return
	}

	// Insert order products
	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			return
		}
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return
	}
	stmt.Close()

	return
}

func (r *postgresRepository) GetOrderByID(ctx context.Context, id string) (*Order, error) {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, created_at, account_id, total_price FROM orders WHERE id = $1",
		id,
	)
	o := &Order{}
	if err := row.Scan(&o.ID, &o.CreatedAt, &o.AccountID, &o.TotalPrice); err != nil {
		return nil, err
	}
	return o, nil
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, created_at, account_id, total_price
    FROM orders
    WHERE account_id = $1`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		o := &Order{}
		if err := rows.Scan(&o.ID, &o.CreatedAt, &o.AccountID, &o.TotalPrice); err == nil {
			orders = append(orders, *o)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
