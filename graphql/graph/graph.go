//go:generate gqlgen -schema ../schema.graphql
package graph

import (
	context "context"
	"time"

	accountProto "github.com/tinrab/spidey/account/pb"
	catalogProto "github.com/tinrab/spidey/catalog/pb"
	"google.golang.org/grpc"
)

type GraphQLServer struct {
	accountConn   *grpc.ClientConn
	accountClient accountProto.AccountServiceClient
	catalogClient catalogProto.CatalogServiceClient
}

func NewGraphQLServer(accountUrl string, catalogUrl string) (*GraphQLServer, error) {
	// Connect to account service
	accountConn, err := grpc.Dial(accountUrl, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	accountClient := accountProto.NewAccountServiceClient(accountConn)

	// Connect to product service
	catalogConn, err := grpc.Dial(catalogUrl, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	catalogClient := catalogProto.NewCatalogServiceClient(catalogConn)

	return &GraphQLServer{accountConn, accountClient, catalogClient}, nil
}

func (s *GraphQLServer) Close() {
	s.accountConn.Close()
}

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

func (s *GraphQLServer) Query_accounts(ctx context.Context, skip *int, take *int, id *string) ([]Account, error) {
	// Get single
	if id != nil {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		r, err := s.accountClient.GetAccount(
			ctx,
			&accountProto.GetAccountRequest{Id: *id},
		)
		if err != nil {
			return nil, err
		}
		return []Account{Account{
			ID:   r.Account.Id,
			Name: r.Account.Name,
		}}, nil
	}

	// Get range
	skipValue := uint64(0)
	takeValue := uint64(100)

	if skip != nil {
		skipValue = uint64(*skip)
	}
	if take != nil {
		takeValue = uint64(*take)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	r, err := s.accountClient.GetAccounts(
		ctx,
		&accountProto.GetAccountsRequest{Skip: skipValue, Take: takeValue},
	)
	if err != nil {
		return nil, err
	}

	accounts := []Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{ID: a.Id, Name: a.Name})
	}

	return accounts, nil
}

func (s *GraphQLServer) Query_products(ctx context.Context, skip *int, take *int, id *string) ([]Product, error) {
	// Get single
	if id != nil {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		r, err := s.catalogClient.GetProduct(
			ctx,
			&catalogProto.GetProductRequest{Id: *id},
		)
		if err != nil {
			return nil, err
		}
		return []Product{Product{
			ID:          r.Product.Id,
			Name:        r.Product.Name,
			Description: r.Product.Description,
			Price:       r.Product.Price,
		}}, nil
	}

	// Get range
	skipValue := uint64(0)
	takeValue := uint64(100)

	if skip != nil {
		skipValue = uint64(*skip)
	}
	if take != nil {
		takeValue = uint64(*take)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	r, err := s.catalogClient.GetProducts(
		ctx,
		&catalogProto.GetProductsRequest{Skip: skipValue, Take: takeValue},
	)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for _, a := range r.Products {
		products = append(
			products,
			Product{
				ID:          a.Id,
				Name:        a.Name,
				Description: a.Description,
				Price:       a.Price,
			},
		)
	}

	return products, nil
}
