package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccount(ctx context.Context, a Account) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &accountService{r}
}

func (s *accountService) PostAccount(ctx context.Context, a Account) (*Account, error) {
	a.ID = ksuid.New().String()
	if err := s.repository.PutAccount(ctx, a); err != nil {
		return nil, err
	}
	newAccount := &Account{ID: a.ID, Name: a.Name}
	return newAccount, nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}
