package usecase

import (
	"clean-architecture/internal/domain/entity"
	"clean-architecture/internal/domain/repository"
)

type CreateOrderInputDTO struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type CreateOrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository repository.OrderRepositoryInterface
}

func NewCreateOrderUseCase(orderRepository repository.OrderRepositoryInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *CreateOrderUseCase) Execute(input CreateOrderInputDTO) (*CreateOrderOutputDTO, error) {
	order := entity.NewOrder(input.Price, input.Tax)

	err := c.OrderRepository.Save(order)
	if err != nil {
		return nil, err
	}

	return &CreateOrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
