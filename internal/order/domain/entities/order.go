package entities

import (
	"time"
)

type OrderEntity struct {
	ID          uint                  `gorm:"primaryKey"`
	CreatedAt   time.Time             `gorm:"default:current_timestamp"`
	TotalAmount float32               `gorm:"default:0"`
	CustomerId  uint                  `gorm:"index"` // No pointer, no FK constraint - references external customer service
	Products    []*OrderProductEntity `gorm:"foreignKey:OrderId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Status      []*OrderStatusEntity  `gorm:"foreignKey:OrderId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (OrderEntity) TableName() string {
	return "order"
}
