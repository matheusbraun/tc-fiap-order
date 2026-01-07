package commands

type UpdateOrderStatusCommand struct {
	OrderId uint
	Status  uint
}

func NewUpdateOrderStatusCommand(orderId uint, status uint) *UpdateOrderStatusCommand {
	return &UpdateOrderStatusCommand{
		OrderId: orderId,
		Status:  status,
	}
}
