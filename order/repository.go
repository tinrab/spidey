package order

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
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
	placeholders := []string{}
	values := []interface{}{}
	for _, p := range o.Products {
		placeholders = append(placeholders, "(?,?,?)")
		values = append(values, p.OrderID, p.ProductID, p.Quantity)
	}

	stmt, err := tx.PrepareContext(
		ctx,
		"INSERT INTO order_products(order_id, product_id, quantity) VALUES "+
			strings.Join(placeholders, ","),
	)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, values...)

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
	return nil, nil
}
