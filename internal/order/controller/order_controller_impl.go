package controller

import (
	"github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/presenter"
	addorder "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/addOrder"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
	getorder "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrder"
	getorderstatus "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrderStatus"
	getorders "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrders"
	updateorderstatus "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/updateOrderStatus"
)

var (
	_ OrderController = (*OrderControllerImpl)(nil)
)

type OrderControllerImpl struct {
	presenter                presenter.OrderPresenter
	addOrderUseCase          addorder.AddOrderUseCase
	getOrderUseCase          getorder.GetOrderUseCase
	getOrdersUseCase         getorders.GetOrdersUseCase
	getOrderStatusUseCase    getorderstatus.GetOrderStatusUseCase
	updateOrderStatusUseCase updateorderstatus.UpdateOrderStatusUseCase
}

func NewOrderControllerImpl(
	presenter presenter.OrderPresenter,
	addOrderUseCase addorder.AddOrderUseCase,
	getOrderUseCase getorder.GetOrderUseCase,
	getOrdersUseCase getorders.GetOrdersUseCase,
	getOrderStatusUseCase getorderstatus.GetOrderStatusUseCase,
	updateOrderStatusUseCase updateorderstatus.UpdateOrderStatusUseCase) *OrderControllerImpl {
	return &OrderControllerImpl{
		presenter:                presenter,
		addOrderUseCase:          addOrderUseCase,
		getOrderUseCase:          getOrderUseCase,
		getOrdersUseCase:         getOrdersUseCase,
		getOrderStatusUseCase:    getOrderStatusUseCase,
		updateOrderStatusUseCase: updateOrderStatusUseCase,
	}
}

func (c *OrderControllerImpl) Add(addOrderRequest *dto.AddOrderDto) (string, error) {
	orderId, err := c.addOrderUseCase.Execute(commands.NewAddOrderCommand(
		*addOrderRequest.CustomerId,
		addOrderRequest.TotalAmount,
		addOrderRequest.Products))
	if err != nil {
		return "", err
	}

	return orderId, nil
}

func (c *OrderControllerImpl) GetOrder(orderId uint) (*dto.GetOrderResponseDto, error) {
	order, err := c.getOrderUseCase.Execute(commands.NewGetOrderCommand(orderId))
	if err != nil {
		return nil, err
	}

	return c.presenter.Present(order), nil
}

func (c *OrderControllerImpl) GetOrders() (*dto.GetOrdersResponseDto, error) {
	orders, err := c.getOrdersUseCase.Execute(commands.NewGetOrdersCommand())
	if err != nil {
		return nil, err
	}

	return c.presenter.PresentOrders(orders), nil
}

func (c *OrderControllerImpl) GetOrderStatus(orderId uint) (*dto.GetOrderStatusResponseDto, error) {
	orderStatus, err := c.getOrderStatusUseCase.Execute(commands.NewGetOrderStatusCommand(orderId))
	if err != nil {
		return nil, err
	}

	return c.presenter.PresentStatus(orderStatus), nil
}

func (c *OrderControllerImpl) UpdateOrderStatus(orderId uint, updateOrderStatusRequest *dto.UpdateOrderStatusRequestDto) error {
	err := c.updateOrderStatusUseCase.Execute(commands.NewUpdateOrderStatusCommand(orderId, updateOrderStatusRequest.Status))
	if err != nil {
		return err
	}

	return nil
}
