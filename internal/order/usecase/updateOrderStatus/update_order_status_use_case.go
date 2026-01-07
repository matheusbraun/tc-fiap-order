package updateorderstatus

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
)

type UpdateOrderStatusUseCase interface {
	Execute(command *commands.UpdateOrderStatusCommand) error
}
