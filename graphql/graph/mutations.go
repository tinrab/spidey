package graph

import (
	context "context"
	time "time"

	accountProto "github.com/tinrab/spidey/account/pb"
	catalogProto "github.com/tinrab/spidey/catalog/pb"
)

func (s *GraphQLServer) Mutation_createAccount(ctx context.Context, name string) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	r, err := s.accountClient.PostAccount(
		ctx,
		&accountProto.PostAccountRequest{Name: name},
	)
	if err != nil {
		return nil, err
	}
	return &Account{ID: r.Account.Id, Name: r.Account.Name}, nil
}

func (s *GraphQLServer) Mutation_createProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	r, err := s.catalogClient.PostProduct(
		ctx,
		&catalogProto.PostProductRequest{Name: name, Description: description, Price: price},
	)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}
