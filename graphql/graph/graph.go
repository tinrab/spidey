//go:generate gqlgen -schema ../schema.graphql
package graph

import (
	context "context"
	"time"

	accountProto "github.com/tinrab/spidey/account/pb"
	"google.golang.org/grpc"
)

type GraphQLServer struct {
	accountConn   *grpc.ClientConn
	accountClient accountProto.AccountServiceClient
}

func NewGraphQLServer(accountUrl string) (*GraphQLServer, error) {
	// Connect to account service
	accountConn, err := grpc.Dial(accountUrl, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	accountClient := accountProto.NewAccountServiceClient(accountConn)

	return &GraphQLServer{accountConn, accountClient}, nil
}

func (s *GraphQLServer) Close() {
	s.accountConn.Close()
}

func (s *GraphQLServer) Mutation_createAccount(ctx context.Context, name string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	r, err := s.accountClient.PostAccount(
		ctx,
		&accountProto.PostAccountRequest{Name: name},
	)
	if err != nil {
		return "", err
	}
	return r.Id, nil
}

func (s *GraphQLServer) Mutation_createProduct(ctx context.Context, name string, description string, price float64) (string, error) {
	return "", nil
}

func (s *GraphQLServer) Query_accounts(ctx context.Context, skip *int, take *int) ([]Account, error) {
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
