package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, o Order) (*Order, error)
	GetOrder(ctx context.Context, id string) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	TotalPrice float64   `json:"total_price"`
	AccountID  string    `json:"account_id"`
	Products   []Product `json:"products"`
}

type Product struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  uint32 `json:"quantity"`
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s orderService) PostOrder(ctx context.Context, o Order) (*Order, error) {
	o.ID = ksuid.New().String()
	if err := s.repository.PutOrder(ctx, o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s orderService) GetOrder(ctx context.Context, id string) (*Order, error) {
	return s.repository.GetOrderByID(ctx, id)
}

func (s orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accountID)
}
