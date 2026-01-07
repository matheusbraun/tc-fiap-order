package entities

type OrderProductEntity struct {
	ID        uint        `gorm:"primaryKey"`
	OrderId   uint        `gorm:"index"`
	ProductId uint        `gorm:"index"` // No FK constraint - references external product service
	Price     float32     `gorm:"not null"`
	Quantity  uint        `gorm:"not null"`
	Order     OrderEntity `gorm:"foreignKey:OrderId;references:ID"`
}

func (OrderProductEntity) TableName() string {
	return "order_product"
}
