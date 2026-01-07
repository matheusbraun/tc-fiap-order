package getorderstatus

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
)

type GetOrderStatusUseCase interface {
	Execute(command *commands.GetOrderStatusCommand) (*entities.OrderStatusEntity, error)
}
