package repositories

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
)

type OrderProductRepository interface {
	AddOrderProduct(orderProduct *entities.OrderProductEntity) error
}
