package entities

import (
	"time"
)

type OrderStatusEntity struct {
	ID            uint        `gorm:"primaryKey"`
	CreatedAt     time.Time   `gorm:"default:current_timestamp"`
	CurrentStatus uint        `gorm:"not null"`
	OrderId       uint        `gorm:"index"`
	Order         OrderEntity `gorm:"foreignKey:OrderId;references:ID"`
}

func (OrderStatusEntity) TableName() string {
	return "order_status"
}
