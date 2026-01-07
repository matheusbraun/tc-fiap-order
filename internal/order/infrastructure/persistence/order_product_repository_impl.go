package secondary

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/repositories"
	"gorm.io/gorm"
)

var (
	_ repositories.OrderProductRepository = (*OrderProductRepositoryImpl)(nil)
)

type OrderProductRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderProductRepositoryImpl(db *gorm.DB) *OrderProductRepositoryImpl {
	return &OrderProductRepositoryImpl{db: db}
}

func (r *OrderProductRepositoryImpl) AddOrderProduct(orderProduct *entities.OrderProductEntity) error {
	if err := r.db.Create(orderProduct).Error; err != nil {
		return err
	}
	return nil
}
