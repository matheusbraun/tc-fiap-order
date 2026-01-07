package repositories

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
)

type OrderStatusRepository interface {
	AddOrderStatus(orderStatus *entities.OrderStatusEntity) error
	GetOrderStatus(orderId uint) (*entities.OrderStatusEntity, error)
}
