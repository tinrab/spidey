//go:generate gqlgen -schema ../schema.graphql -typemap ../types.json
package graph

import (
	"github.com/tinrab/spidey/account"
	"github.com/tinrab/spidey/catalog"
	"github.com/tinrab/spidey/order"
)

type GraphQLServer struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQLServer(accountUrl, catalogURL, orderURL string) (*GraphQLServer, error) {
	// Connect to account service
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	// Connect to product service
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	// Connect to order service
	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}

	return &GraphQLServer{
		accountClient,
		catalogClient,
		orderClient,
	}, nil
}
