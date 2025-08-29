package database

import (
	"database/sql"

	"clean-architecture/internal/domain/entity"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.DB.Prepare("INSERT INTO orders (id, price, tax, final_price, created_at) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice, order.CreatedAt)
	return err
}

func (r *OrderRepository) FindAll() ([]*entity.Order, error) {
	rows, err := r.DB.Query("SELECT id, price, tax, final_price, created_at FROM orders ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		order := &entity.Order{}
		err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) FindByID(id string) (*entity.Order, error) {
	order := &entity.Order{}
	err := r.DB.QueryRow("SELECT id, price, tax, final_price, created_at FROM orders WHERE id = $1", id).
		Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice, &order.CreatedAt)
	if err != nil {
		return nil, err
	}
	return order, nil
}
