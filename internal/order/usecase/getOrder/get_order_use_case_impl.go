package getorder

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/repositories"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
)

var (
	_ GetOrderUseCase = (*GetOrderUseCaseImpl)(nil)
)

type GetOrderUseCaseImpl struct {
	orderRepository repositories.OrderRepository
}

func NewGetOrderUseCaseImpl(orderRepository repositories.OrderRepository) *GetOrderUseCaseImpl {
	return &GetOrderUseCaseImpl{orderRepository: orderRepository}
}

func (u *GetOrderUseCaseImpl) Execute(command *commands.GetOrderCommand) (*entities.OrderEntity, error) {
	order, err := u.orderRepository.GetOrder(command.OrderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}
