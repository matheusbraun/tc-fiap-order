package dto

import (
	"time"
)

type CustomerDto struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	CPF   uint   `json:"cpf"`
}

type GetOrderResponseDto struct {
	ID          uint                         `json:"id"`
	CreatedAt   time.Time                    `json:"created_at"`
	TotalAmount float32                      `json:"total_amount"`
	CustomerId  uint                         `json:"customer_id,omitempty"`
	Customer    *CustomerDto                 `json:"customer,omitempty"`
	Products    []*OrderProductDto           `json:"products"`
	Status      []*GetOrderStatusResponseDto `json:"status"`
}

type OrderProductDto struct {
	ProductId   uint    `json:"product_id"`
	Price       float32 `json:"price"`
	Quantity    uint    `json:"quantity"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Category    int     `json:"category,omitempty"`
	ImageLink   string  `json:"image_link,omitempty"`
}
