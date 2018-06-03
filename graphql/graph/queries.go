package graph

import (
	context "context"
	"log"
	time "time"
)

func (s *GraphQLServer) Account_orders(ctx context.Context, obj *Account) ([]Order, error) {
	return nil, nil
}

func (s *GraphQLServer) Order_products(ctx context.Context, obj *Order) ([]OrderedProduct, error) {
	return nil, nil
}

func (s *GraphQLServer) Query_accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get single
	if id != nil {
		r, err := s.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return []Account{Account{
			ID:   r.ID,
			Name: r.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	r, err := s.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	accounts := []Account{}
	for _, a := range r {
		accounts = append(accounts, Account{ID: a.ID, Name: a.Name})
	}

	return accounts, nil
}

func (s *GraphQLServer) Query_products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get single
	if id != nil {
		r, err := s.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []Product{Product{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	r, err := s.catalogClient.GetProducts(ctx, skip, take, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []Product{}
	for _, a := range r {
		products = append(products,
			Product{
				ID:          a.ID,
				Name:        a.Name,
				Description: a.Description,
				Price:       a.Price,
			},
		)
	}

	return products, nil
}

func (p PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}
	if p.Take != nil {
		takeValue = uint64(*p.Take)
	}
	return skipValue, takeValue
}
