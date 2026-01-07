package addorder

import (
	"fmt"

	"github.com/viniciuscluna/tc-fiap-50/internal/infrastructure/clients"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/repositories"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
)

var (
	_ AddOrderUseCase = (*AddOrderUseCaseImpl)(nil)
)

type AddOrderUseCaseImpl struct {
	orderRepository        repositories.OrderRepository
	orderProductRepository repositories.OrderProductRepository
	orderStatusRepository  repositories.OrderStatusRepository
	customerClient         clients.CustomerClient
	productClient          clients.ProductClient
}

func NewAddOrderUseCaseImpl(
	orderRepository repositories.OrderRepository,
	orderProductRepository repositories.OrderProductRepository,
	orderStatusRepository repositories.OrderStatusRepository,
	customerClient clients.CustomerClient,
	productClient clients.ProductClient) *AddOrderUseCaseImpl {
	return &AddOrderUseCaseImpl{
		orderRepository:        orderRepository,
		orderProductRepository: orderProductRepository,
		orderStatusRepository:  orderStatusRepository,
		customerClient:         customerClient,
		productClient:          productClient,
	}
}

func (u *AddOrderUseCaseImpl) Execute(command *commands.AddOrderCommand) (string, error) {
	// Extract product IDs and validate all exist
	productIDs := make([]uint, len(command.Products))
	for i, p := range command.Products {
		productIDs[i] = p.ProductId
	}

	// Create order
	orderResult, err := u.orderRepository.AddOrder(&entities.OrderEntity{
		CustomerId:  command.CustomerId,
		TotalAmount: command.TotalAmount,
	})
	if err != nil {
		return "", err
	}

	// Add order products
	for _, orderProductDto := range command.Products {
		orderProductEntity := &entities.OrderProductEntity{
			OrderId:   orderResult.ID,
			ProductId: orderProductDto.ProductId,
			Price:     orderProductDto.Price,
			Quantity:  orderProductDto.Quantity,
		}
		err := u.orderProductRepository.AddOrderProduct(orderProductEntity)
		if err != nil {
			return "", err
		}
	}

	// Create initial order status
	orderStatusEntity := &entities.OrderStatusEntity{
		OrderId:       orderResult.ID,
		CurrentStatus: 1,
	}
	err = u.orderStatusRepository.AddOrderStatus(orderStatusEntity)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", orderResult.ID), nil
}
