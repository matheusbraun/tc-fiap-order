package getorders

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/repositories"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
)

var (
	_ GetOrdersUseCase = (*GetOrdersUseCaseImpl)(nil)
)

type GetOrdersUseCaseImpl struct {
	orderRepository repositories.OrderRepository
}

func NewGetOrdersUseCaseImpl(orderRepository repositories.OrderRepository) *GetOrdersUseCaseImpl {
	return &GetOrdersUseCaseImpl{orderRepository: orderRepository}
}

func (u *GetOrdersUseCaseImpl) Execute(command *commands.GetOrdersCommand) ([]*entities.OrderEntity, error) {
	orders, err := u.orderRepository.GetOrders()
	if err != nil {
		return nil, err
	}

	return orders, nil
}
