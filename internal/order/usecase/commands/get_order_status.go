package commands

type GetOrderStatusCommand struct {
	OrderId uint
}

func NewGetOrderStatusCommand(orderId uint) *GetOrderStatusCommand {
	return &GetOrderStatusCommand{
		OrderId: orderId,
	}
}
