package updateorderstatus

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/repositories"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
)

var (
	_ UpdateOrderStatusUseCase = (*UpdateOrderStatusUseCaseImpl)(nil)
)

type UpdateOrderStatusUseCaseImpl struct {
	orderStatusRepository repositories.OrderStatusRepository
}

func NewUpdateOrderStatusUseCaseImpl(orderStatusRepository repositories.OrderStatusRepository) *UpdateOrderStatusUseCaseImpl {
	return &UpdateOrderStatusUseCaseImpl{
		orderStatusRepository: orderStatusRepository,
	}
}

func (u *UpdateOrderStatusUseCaseImpl) Execute(command *commands.UpdateOrderStatusCommand) error {
	err := u.orderStatusRepository.AddOrderStatus(&entities.OrderStatusEntity{
		OrderId:       command.OrderId,
		CurrentStatus: command.Status,
	})
	if err != nil {
		return err
	}

	return nil
}
