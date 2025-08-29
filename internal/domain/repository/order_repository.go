package repository

import "clean-architecture/internal/domain/entity"

type OrderRepositoryInterface interface {
	Save(order *entity.Order) error
	FindAll() ([]*entity.Order, error)
	FindByID(id string) (*entity.Order, error)
}
