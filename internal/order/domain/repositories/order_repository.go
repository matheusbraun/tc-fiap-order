package repositories

import "github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"

type OrderRepository interface {
	AddOrder(order *entities.OrderEntity) (*entities.OrderEntity, error)
	GetOrder(orderId uint) (*entities.OrderEntity, error)
	GetOrders() ([]*entities.OrderEntity, error)
}
