package usecase

import (
	"clean-architecture/internal/domain/repository"
)

type ListOrdersOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrdersUseCase struct {
	OrderRepository repository.OrderRepositoryInterface
}

func NewListOrdersUseCase(orderRepository repository.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (l *ListOrdersUseCase) Execute() ([]ListOrdersOutputDTO, error) {
	orders, err := l.OrderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var result []ListOrdersOutputDTO
	for _, order := range orders {
		result = append(result, ListOrdersOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return result, nil
}
