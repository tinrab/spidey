//go:generate protoc ./order.proto --go_out=plugins=grpc:./pb
package order

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/tinrab/spidey/catalog"
	"github.com/tinrab/spidey/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, catalogUrl string, port int) error {
	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{s, catalogClient})
	reflection.Register(serv)

	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(
	ctx context.Context,
	r *pb.PostOrderRequest,
) (*pb.PostOrderResponse, error) {
	// Get ordered products
	productIDs := []string{}
	for _, p := range r.Products {
		productIDs = append(productIDs, p.ProductId)
	}
	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs)
	if err != nil {
		return nil, err
	}

	products := []OrderedProduct{}
	totalPrice := 0.0
	for _, p := range orderedProducts {
		// Set product if it exists
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

	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, o.CreatedAt)
	op.CreatedAt = buf.Bytes()

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
