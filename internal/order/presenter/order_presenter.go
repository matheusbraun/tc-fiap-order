package presenter

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"
)

type OrderPresenter interface {
	Present(order *entities.OrderEntity) *dto.GetOrderResponseDto
	PresentOrders(orders []*entities.OrderEntity) *dto.GetOrdersResponseDto
	PresentProducts(orderProducts []*entities.OrderProductEntity) []*dto.OrderProductDto
	PresentStatus(orderStatus *entities.OrderStatusEntity) *dto.GetOrderStatusResponseDto
	PresentMultipleStatus(orderStatus []*entities.OrderStatusEntity) []*dto.GetOrderStatusResponseDto
}
