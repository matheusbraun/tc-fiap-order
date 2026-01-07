package commands

import "github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"

type AddOrderCommand struct {
	CustomerId  uint
	TotalAmount float32
	Products    []*dto.AddOrderProductDto
}

func NewAddOrderCommand(customerId uint, totalAmount float32, products []*dto.AddOrderProductDto) *AddOrderCommand {
	return &AddOrderCommand{
		CustomerId:  customerId,
		TotalAmount: totalAmount,
		Products:    products,
	}
}
