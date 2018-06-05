//go:generate protoc ./order.proto --go_out=plugins=grpc:./pb
package order

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/tinrab/spidey/account"
	"github.com/tinrab/spidey/catalog"
	"github.com/tinrab/spidey/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		s,
		accountClient,
		catalogClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(
	ctx context.Context,
	r *pb.PostOrderRequest,
) (*pb.PostOrderResponse, error) {
	// Check if account exists
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Get ordered products
	productIDs := []string{}
	for _, p := range r.Products {
		productIDs = append(productIDs, p.ProductId)
	}
	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []OrderedProduct{}
	totalPrice := 0.0
	for _, p := range orderedProducts {
		// Include product if it exists
		quantity := uint32(0)
		for _, rp := range r.Products {
			if rp.ProductId == p.ID {
				quantity = rp.Quantity
				break
			}
		}
		if quantity != 0 {
			products = append(products, OrderedProduct{
				ID:       p.ID,
				Quantity: quantity,
			})
		}
		// Calculate total price
		totalPrice += p.Price
	}

	o, err := s.service.PostOrder(ctx, r.AccountId, totalPrice, products)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.PostOrderResponse{
		Order: s.encodeOrder(*o),
	}, nil
}

func (s *grpcServer) GetOrder(
	ctx context.Context,
	r *pb.GetOrderRequest,
) (*pb.GetOrderResponse, error) {
	o, err := s.service.GetOrder(ctx, r.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.GetOrderResponse{
		Order: s.encodeOrder(*o),
	}, nil
}

func (s *grpcServer) GetOrdersForAccount(
	ctx context.Context,
	r *pb.GetOrdersForAccountRequest,
) (*pb.GetOrdersForAccountResponse, error) {
	res, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	orders := []*pb.Order{}
	for _, o := range res {
		orders = append(orders, s.encodeOrder(o))
	}
	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}

func (s *grpcServer) encodeOrder(o Order) *pb.Order {
	op := &pb.Order{
		AccountId:  o.AccountID,
		Id:         o.ID,
		TotalPrice: o.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

	for _, p := range o.Products {
		op.Products = append(op.Products, &pb.Order_OrderProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		})
	}
	return op
}

func (s *grpcServer) decodeOrder(r *pb.PostOrderRequest) Order {
	o := Order{
		AccountID: r.AccountId,
		Products:  []OrderedProduct{},
	}

	for _, p := range r.Products {
		o.Products = append(o.Products, OrderedProduct{
			ID:       p.ProductId,
			Quantity: p.Quantity,
		})
	}
	return o
}
