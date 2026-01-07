package dto

type AddOrderDto struct {
	CustomerId  *uint                 `json:"customerId" example:"1"`
	TotalAmount float32               `json:"totalAmount" example:"34.99"`
	Products    []*AddOrderProductDto `json:"products"`
}

type AddOrderProductDto struct {
	ProductId uint    `json:"productId" example:"1"`
	Quantity  uint    `json:"quantity" example:"1"`
	Price     float32 `json:"price" example:"34.99"`
}
