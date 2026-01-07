package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/controller"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"
	"github.com/viniciuscluna/tc-fiap-50/tests/mocks"
)

func TestOrderController_Add(t *testing.T) {
	t.Run("deve adicionar pedido com sucesso", func(t *testing.T) {
		// Setup
		addOrderUseCase := mocks.NewMockAddOrderUseCase(t)
		presenter := mocks.NewMockOrderPresenter(t)
		getOrderUseCase := mocks.NewMockGetOrderUseCase(t)
		getOrdersUseCase := mocks.NewMockGetOrdersUseCase(t)
		getOrderStatusUseCase := mocks.NewMockGetOrderStatusUseCase(t)
		updateOrderStatusUseCase := mocks.NewMockUpdateOrderStatusUseCase(t)

		ctrl := controller.NewOrderControllerImpl(
			presenter,
			addOrderUseCase,
			getOrderUseCase,
			getOrdersUseCase,
			getOrderStatusUseCase,
			updateOrderStatusUseCase,
		)

		customerId := uint(1)
		addOrderDto := &dto.AddOrderDto{
			CustomerId:  &customerId,
			TotalAmount: 100.0,
			Products: []*dto.AddOrderProductDto{
				{ProductId: 101, Quantity: 2},
			},
		}

		addOrderUseCase.On("Execute", mock.MatchedBy(func(cmd interface{}) bool {
			return true
		})).Return("1", nil)

		// Execute
		result, err := ctrl.Add(addOrderDto)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "1", result)
		addOrderUseCase.AssertExpectations(t)
	})

	t.Run("deve retornar erro ao falhar no use case", func(t *testing.T) {
		// Setup
		addOrderUseCase := mocks.NewMockAddOrderUseCase(t)
		presenter := mocks.NewMockOrderPresenter(t)
		getOrderUseCase := mocks.NewMockGetOrderUseCase(t)
		getOrdersUseCase := mocks.NewMockGetOrdersUseCase(t)
		getOrderStatusUseCase := mocks.NewMockGetOrderStatusUseCase(t)
		updateOrderStatusUseCase := mocks.NewMockUpdateOrderStatusUseCase(t)

		ctrl := controller.NewOrderControllerImpl(
			presenter,
			addOrderUseCase,
			getOrderUseCase,
			getOrdersUseCase,
			getOrderStatusUseCase,
			updateOrderStatusUseCase,
		)

		customerId := uint(1)
		addOrderDto := &dto.AddOrderDto{
			CustomerId:  &customerId,
			TotalAmount: 100.0,
			Products:    []*dto.AddOrderProductDto{},
		}

		addOrderUseCase.On("Execute", mock.Anything).Return("", assert.AnError)

		// Execute
		result, err := ctrl.Add(addOrderDto)

		// Assert
		require.Error(t, err)
		assert.Empty(t, result)
		addOrderUseCase.AssertExpectations(t)
	})
}

func TestOrderController_GetOrder(t *testing.T) {
	t.Run("deve obter pedido com sucesso", func(t *testing.T) {
		// Setup
		addOrderUseCase := mocks.NewMockAddOrderUseCase(t)
		presenter := mocks.NewMockOrderPresenter(t)
		getOrderUseCase := mocks.NewMockGetOrderUseCase(t)
		getOrdersUseCase := mocks.NewMockGetOrdersUseCase(t)
		getOrderStatusUseCase := mocks.NewMockGetOrderStatusUseCase(t)
		updateOrderStatusUseCase := mocks.NewMockUpdateOrderStatusUseCase(t)

		ctrl := controller.NewOrderControllerImpl(
			presenter,
			addOrderUseCase,
			getOrderUseCase,
			getOrdersUseCase,
			getOrderStatusUseCase,
			updateOrderStatusUseCase,
		)

		orderEntity := &entities.OrderEntity{
			ID:          1,
			CustomerId:  1,
			TotalAmount: 100.0,
		}

		expectedDto := &dto.GetOrderResponseDto{
			ID:          1,
			TotalAmount: 100.0,
		}

		getOrderUseCase.On("Execute", mock.Anything).Return(orderEntity, nil)
		presenter.On("Present", orderEntity).Return(expectedDto)

		// Execute
		result, err := ctrl.GetOrder(1)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedDto, result)
		getOrderUseCase.AssertExpectations(t)
		presenter.AssertExpectations(t)
	})

	t.Run("deve retornar erro ao falhar no use case", func(t *testing.T) {
		// Setup
		addOrderUseCase := mocks.NewMockAddOrderUseCase(t)
		presenter := mocks.NewMockOrderPresenter(t)
		getOrderUseCase := mocks.NewMockGetOrderUseCase(t)
		getOrdersUseCase := mocks.NewMockGetOrdersUseCase(t)
		getOrderStatusUseCase := mocks.NewMockGetOrderStatusUseCase(t)
		updateOrderStatusUseCase := mocks.NewMockUpdateOrderStatusUseCase(t)

		ctrl := controller.NewOrderControllerImpl(
			presenter,
			addOrderUseCase,
			getOrderUseCase,
			getOrdersUseCase,
			getOrderStatusUseCase,
			updateOrderStatusUseCase,
		)

		getOrderUseCase.On("Execute", mock.Anything).Return(nil, assert.AnError)

		// Execute
		result, err := ctrl.GetOrder(999)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		getOrderUseCase.AssertExpectations(t)
	})
}

func TestOrderController_GetOrders(t *testing.T) {
	t.Run("deve listar pedidos com sucesso", func(t *testing.T) {
		// Setup
		addOrderUseCase := mocks.NewMockAddOrderUseCase(t)
		presenter := mocks.NewMockOrderPresenter(t)
		getOrderUseCase := mocks.NewMockGetOrderUseCase(t)
		getOrdersUseCase := mocks.NewMockGetOrdersUseCase(t)
		getOrderStatusUseCase := mocks.NewMockGetOrderStatusUseCase(t)
		updateOrderStatusUseCase := mocks.NewMockUpdateOrderStatusUseCase(t)

		ctrl := controller.NewOrderControllerImpl(
			presenter,
			addOrderUseCase,
			getOrderUseCase,
			getOrdersUseCase,
			getOrderStatusUseCase,
			updateOrderStatusUseCase,
		)

		orders := []*entities.OrderEntity{
			{ID: 1, CustomerId: 1, TotalAmount: 100.0},
			{ID: 2, CustomerId: 2, TotalAmount: 200.0},
		}

		expectedDto := &dto.GetOrdersResponseDto{
			Orders: []*dto.GetOrderResponseDto{
				{ID: 1, TotalAmount: 100.0},
				{ID: 2, TotalAmount: 200.0},
			},
		}

		getOrdersUseCase.On("Execute", mock.Anything).Return(orders, nil)
		presenter.On("PresentOrders", orders).Return(expectedDto)

		// Execute
		result, err := ctrl.GetOrders()

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedDto, result)
		getOrdersUseCase.AssertExpectations(t)
		presenter.AssertExpectations(t)
	})
}

func TestOrderController_GetOrderStatus(t *testing.T) {
	t.Run("deve obter status do pedido com sucesso", func(t *testing.T) {
		// Setup
		addOrderUseCase := mocks.NewMockAddOrderUseCase(t)
		presenter := mocks.NewMockOrderPresenter(t)
		getOrderUseCase := mocks.NewMockGetOrderUseCase(t)
		getOrdersUseCase := mocks.NewMockGetOrdersUseCase(t)
		getOrderStatusUseCase := mocks.NewMockGetOrderStatusUseCase(t)
		updateOrderStatusUseCase := mocks.NewMockUpdateOrderStatusUseCase(t)

		ctrl := controller.NewOrderControllerImpl(
			presenter,
			addOrderUseCase,
			getOrderUseCase,
			getOrdersUseCase,
			getOrderStatusUseCase,
			updateOrderStatusUseCase,
		)

		statusEntity := &entities.OrderStatusEntity{
			ID:            1,
			OrderId:       1,
			CurrentStatus: 1,
		}

		expectedDto := &dto.GetOrderStatusResponseDto{
			ID:                       1,
			OrderId:                  1,
			CurrentStatus:            1,
			CurrentStatusDescription: "Recebido",
		}

		getOrderStatusUseCase.On("Execute", mock.Anything).Return(statusEntity, nil)
		presenter.On("PresentStatus", statusEntity).Return(expectedDto)

		// Execute
		result, err := ctrl.GetOrderStatus(1)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedDto, result)
		getOrderStatusUseCase.AssertExpectations(t)
		presenter.AssertExpectations(t)
	})
}

func TestOrderController_UpdateOrderStatus(t *testing.T) {
	t.Run("deve atualizar status do pedido com sucesso", func(t *testing.T) {
		// Setup
		addOrderUseCase := mocks.NewMockAddOrderUseCase(t)
		presenter := mocks.NewMockOrderPresenter(t)
		getOrderUseCase := mocks.NewMockGetOrderUseCase(t)
		getOrdersUseCase := mocks.NewMockGetOrdersUseCase(t)
		getOrderStatusUseCase := mocks.NewMockGetOrderStatusUseCase(t)
		updateOrderStatusUseCase := mocks.NewMockUpdateOrderStatusUseCase(t)

		ctrl := controller.NewOrderControllerImpl(
			presenter,
			addOrderUseCase,
			getOrderUseCase,
			getOrdersUseCase,
			getOrderStatusUseCase,
			updateOrderStatusUseCase,
		)

		updateDto := &dto.UpdateOrderStatusRequestDto{
			Status: 2, // Em Preparação
		}

		updateOrderStatusUseCase.On("Execute", mock.Anything).Return(nil)

		// Execute
		err := ctrl.UpdateOrderStatus(1, updateDto)

		// Assert
		require.NoError(t, err)
		updateOrderStatusUseCase.AssertExpectations(t)
	})

	t.Run("deve retornar erro ao falhar no use case", func(t *testing.T) {
		// Setup
		addOrderUseCase := mocks.NewMockAddOrderUseCase(t)
		presenter := mocks.NewMockOrderPresenter(t)
		getOrderUseCase := mocks.NewMockGetOrderUseCase(t)
		getOrdersUseCase := mocks.NewMockGetOrdersUseCase(t)
		getOrderStatusUseCase := mocks.NewMockGetOrderStatusUseCase(t)
		updateOrderStatusUseCase := mocks.NewMockUpdateOrderStatusUseCase(t)

		ctrl := controller.NewOrderControllerImpl(
			presenter,
			addOrderUseCase,
			getOrderUseCase,
			getOrdersUseCase,
			getOrderStatusUseCase,
			updateOrderStatusUseCase,
		)

		updateDto := &dto.UpdateOrderStatusRequestDto{
			Status: 999, // Status Inválido
		}

		updateOrderStatusUseCase.On("Execute", mock.Anything).Return(assert.AnError)

		// Execute
		err := ctrl.UpdateOrderStatus(999, updateDto)

		// Assert
		require.Error(t, err)
		updateOrderStatusUseCase.AssertExpectations(t)
	})
}
