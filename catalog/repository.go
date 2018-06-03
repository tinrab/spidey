package catalog

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
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

func (r *postgresRepository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO products(id, name, description, price) VALUES($1, $2, $3, $4)",
		p.ID,
		p.Name,
		p.Description,
		p.Price,
	)
	return err
}

func (r *postgresRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, name, description, price FROM products WHERE id = $1",
		id,
	)
	p := &Product{}
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *postgresRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
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
		p := &Product{}
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err == nil {
			products = append(products, *p)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *postgresRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	// Make ID list
	idList := strings.Builder{}
	for i, id := range ids {
		idList.WriteString("\"" + id + "\"")
		if i < len(ids)-1 {
			idList.WriteString(",")
		}
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name, description, price
     FROM products
     WHERE id IN($1)
     ORDER BY id DESC`,
		idList.String(),
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		p := &Product{}
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err == nil {
			products = append(products, *p)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
