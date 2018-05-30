//go:generate protoc ./order.proto --go_out=plugins=grpc:./pb
package order

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/tinrab/spidey/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(
	ctx context.Context,
	r *pb.PostOrderRequest,
) (*pb.PostOrderResponse, error) {
	o, err := s.service.PostOrder(
		ctx,
		s.decodeOrder(r),
	)
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
			OrderId:   p.OrderID,
			ProductId: p.ProductID,
			Quantity:  p.Quantity,
		})
	}
	return op
}

func (s *grpcServer) decodeOrder(r *pb.PostOrderRequest) Order {
	o := Order{
		AccountID:  r.AccountId,
		TotalPrice: r.TotalPrice,
		Products:   []Product{},
	}

	buf := bytes.NewReader(r.CreatedAt)
	binary.Read(buf, binary.LittleEndian, &o.CreatedAt)

	for _, p := range r.Products {
		o.Products = append(o.Products, Product{
			ProductID: p.ProductId,
			Quantity:  p.Quantity,
		})
	}
	return o
}
