package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         string    `json:"id" db:"id"`
	Price      float64   `json:"price" db:"price"`
	Tax        float64   `json:"tax" db:"tax"`
	FinalPrice float64   `json:"final_price" db:"final_price"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

func NewOrder(price, tax float64) *Order {
	return &Order{
		ID:         uuid.New().String(),
		Price:      price,
		Tax:        tax,
		FinalPrice: price + tax,
		CreatedAt:  time.Now(),
	}
}

func (o *Order) CalculateFinalPrice() {
	o.FinalPrice = o.Price + o.Tax
}
