package getorder

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
)

type GetOrderUseCase interface {
	Execute(command *commands.GetOrderCommand) (*entities.OrderEntity, error)
}
