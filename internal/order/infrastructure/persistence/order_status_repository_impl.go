package secondary

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/repositories"
	"gorm.io/gorm"
)

var (
	_ repositories.OrderStatusRepository = (*OrderStatusRepositoryImpl)(nil)
)

type OrderStatusRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderStatusRepositoryImpl(db *gorm.DB) *OrderStatusRepositoryImpl {
	return &OrderStatusRepositoryImpl{db: db}
}

func (r *OrderStatusRepositoryImpl) AddOrderStatus(orderStatus *entities.OrderStatusEntity) error {
	if err := r.db.Create(orderStatus).Error; err != nil {
		return err
	}
	return nil
}

func (r *OrderStatusRepositoryImpl) GetOrderStatus(orderId uint) (*entities.OrderStatusEntity, error) {
	orderStatus := &entities.OrderStatusEntity{}
	if err := r.db.Where("order_id = ?", orderId).Last(orderStatus).Error; err != nil {
		return nil, err
	}
	return orderStatus, nil
}
