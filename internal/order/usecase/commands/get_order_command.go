package commands

type GetOrderCommand struct {
	OrderId uint
}

func NewGetOrderCommand(orderId uint) *GetOrderCommand {
	return &GetOrderCommand{
		OrderId: orderId,
	}
}
