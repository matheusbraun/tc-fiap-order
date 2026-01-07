package clients

import "context"

type CustomerDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	CPF   uint   `json:"cpf"`
	Email string `json:"email"`
}

type CustomerClient interface {
	GetCustomer(ctx context.Context, customerID uint) (*CustomerDTO, error)
}
