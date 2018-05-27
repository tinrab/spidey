package catalog

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	Ping() error
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *PostgresRepository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO products(id, name, description, price)",
		p.ID,
		p.Name,
		p.Description,
		p.Price,
	)
	return err
}

func (r *PostgresRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, description, price FROM products WHERE id = $1", id)
	p := &Product{}
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostgresRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name, description, price FROM products ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip,
		take,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		p := Product{}
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err == nil {
			products = append(products, p)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
