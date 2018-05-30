//go:generate gqlgen -schema ../schema.graphql
package graph

import (
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
