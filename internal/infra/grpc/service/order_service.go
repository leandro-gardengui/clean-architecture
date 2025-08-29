package service

import (
	"context"

	"clean-architecture/internal/infra/grpc/pb"
	"clean-architecture/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewOrderService(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.CreateOrderInputDTO{
		Price: in.GetPrice(),
		Tax:   in.GetTax(),
	}

	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	output, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var orders []*pb.Order
	for _, dto := range output {
		orders = append(orders, &pb.Order{
			Id:         dto.ID,
			Price:      dto.Price,
			Tax:        dto.Tax,
			FinalPrice: dto.FinalPrice,
		})
	}

	return &pb.ListOrdersResponse{
		Orders: orders,
	}, nil
}
