package getorderstatus

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/repositories"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
)

var (
	_ GetOrderStatusUseCase = (*GetOrderStatusUseCaseImpl)(nil)
)

type GetOrderStatusUseCaseImpl struct {
	orderStatusRepository repositories.OrderStatusRepository
}

func NewGetOrderStatusUseCaseImpl(orderStatusRepository repositories.OrderStatusRepository) *GetOrderStatusUseCaseImpl {
	return &GetOrderStatusUseCaseImpl{orderStatusRepository: orderStatusRepository}
}

func (u *GetOrderStatusUseCaseImpl) Execute(command *commands.GetOrderStatusCommand) (*entities.OrderStatusEntity, error) {
	orderStatus, err := u.orderStatusRepository.GetOrderStatus(command.OrderId)
	if err != nil {
		return nil, err
	}

	return orderStatus, nil
}
