package addorder

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
)

type AddOrderUseCase interface {
	Execute(command *commands.AddOrderCommand) (string, error)
}
