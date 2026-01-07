package getorders

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
)

type GetOrdersUseCase interface {
	Execute(command *commands.GetOrdersCommand) ([]*entities.OrderEntity, error)
}
