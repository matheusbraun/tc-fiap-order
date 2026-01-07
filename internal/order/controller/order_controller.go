package controller

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"
)

type OrderController interface {
	Add(addOrderRequest *dto.AddOrderDto) (string, error)
	GetOrder(orderId uint) (*dto.GetOrderResponseDto, error)
	GetOrders() (*dto.GetOrdersResponseDto, error)
	GetOrderStatus(orderId uint) (*dto.GetOrderStatusResponseDto, error)
	UpdateOrderStatus(orderId uint, updateOrderStatusRequest *dto.UpdateOrderStatusRequestDto) error
}
