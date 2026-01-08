package secondary

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/repositories"
	"gorm.io/gorm"
)

var (
	_ repositories.OrderRepository = (*OrderRepositoryImpl)(nil)
)

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepositoryImpl(db *gorm.DB) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{db: db}
}

func (r *OrderRepositoryImpl) AddOrder(order *entities.OrderEntity) (*entities.OrderEntity, error) {
	if err := r.db.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepositoryImpl) GetOrder(orderId uint) (*entities.OrderEntity, error) {
	order := &entities.OrderEntity{}
	if err := r.db.
		Preload("Products").
		Preload("Status", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Where("id = ?", orderId).
		First(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepositoryImpl) GetOrders() ([]*entities.OrderEntity, error) {
	var orders []*entities.OrderEntity
	if err := r.db.
		Preload("Products").
		Preload("Status", func(db *gorm.DB) *gorm.DB {
			return db.Order("current_status DESC")
		}).
		Where("id NOT IN (SELECT order_id FROM order_status WHERE current_status = 4)").
		Order("created_at ASC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
