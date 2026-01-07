package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/viniciuscluna/tc-fiap-50/internal/infrastructure/clients"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"
	addorder "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/addOrder"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
	getorder "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrder"
	getorderstatus "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrderStatus"
	getorders "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrders"
	updateorderstatus "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/updateOrderStatus"
	"github.com/viniciuscluna/tc-fiap-50/tests/mocks"
)

func TestAddOrderUseCase_Execute(t *testing.T) {
	t.Run("deve criar pedido com sucesso", func(t *testing.T) {
		// Arrange
		orderRepo := mocks.NewMockOrderRepository(t)
		orderProductRepo := mocks.NewMockOrderProductRepository(t)
		orderStatusRepo := mocks.NewMockOrderStatusRepository(t)
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		customerClient.On("ValidateCustomer", mock.Anything, uint(1)).Return(true, nil)

		productClient.On("GetProducts", mock.Anything, []uint{101, 102}).
			Return([]*clients.ProductDTO{
				{ID: 101, Name: "Pizza", Price: 25.50},
				{ID: 102, Name: "Suco", Price: 8.00},
			}, nil)

		orderRepo.On("AddOrder", mock.MatchedBy(func(order *entities.OrderEntity) bool {
			return order.CustomerId == 1 && order.TotalAmount == 59.00
		})).Return(&entities.OrderEntity{ID: 1, CustomerId: 1, TotalAmount: 59.00}, nil)

		orderProductRepo.On("AddOrderProduct", mock.Anything).Return(nil).Times(2)
		orderStatusRepo.On("AddOrderStatus", mock.MatchedBy(func(status *entities.OrderStatusEntity) bool {
			return status.OrderId == 1 && status.CurrentStatus == 1
		})).Return(nil)

		useCase := addorder.NewAddOrderUseCaseImpl(
			orderRepo,
			orderProductRepo,
			orderStatusRepo,
			customerClient,
			productClient,
		)

		products := []*dto.AddOrderProductDto{
			{ProductId: 101, Quantity: 2, Price: 25.50},
			{ProductId: 102, Quantity: 1, Price: 8.00},
		}
		command := commands.NewAddOrderCommand(1, 59.00, products)

		// Act
		result, err := useCase.Execute(command)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "1", result)
		orderRepo.AssertExpectations(t)
		orderProductRepo.AssertExpectations(t)
		orderStatusRepo.AssertExpectations(t)
		customerClient.AssertExpectations(t)
		productClient.AssertExpectations(t)
	})

	t.Run("deve falhar quando cliente não existe", func(t *testing.T) {
		// Arrange
		orderRepo := mocks.NewMockOrderRepository(t)
		orderProductRepo := mocks.NewMockOrderProductRepository(t)
		orderStatusRepo := mocks.NewMockOrderStatusRepository(t)
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		customerClient.On("ValidateCustomer", mock.Anything, uint(999)).Return(false, nil)

		useCase := addorder.NewAddOrderUseCaseImpl(
			orderRepo,
			orderProductRepo,
			orderStatusRepo,
			customerClient,
			productClient,
		)

		products := []*dto.AddOrderProductDto{
			{ProductId: 101, Quantity: 1, Price: 25.50},
		}
		command := commands.NewAddOrderCommand(999, 25.50, products)

		// Act
		_, err := useCase.Execute(command)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "customer 999 not found")
	})

	t.Run("deve falhar quando produtos não existem", func(t *testing.T) {
		// Arrange
		orderRepo := mocks.NewMockOrderRepository(t)
		orderProductRepo := mocks.NewMockOrderProductRepository(t)
		orderStatusRepo := mocks.NewMockOrderStatusRepository(t)
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		customerClient.On("ValidateCustomer", mock.Anything, uint(1)).Return(true, nil)
		productClient.On("GetProducts", mock.Anything, []uint{999}).
			Return([]*clients.ProductDTO{}, nil) // Lista vazia = produtos não encontrados

		useCase := addorder.NewAddOrderUseCaseImpl(
			orderRepo,
			orderProductRepo,
			orderStatusRepo,
			customerClient,
			productClient,
		)

		products := []*dto.AddOrderProductDto{
			{ProductId: 999, Quantity: 1, Price: 25.50},
		}
		command := commands.NewAddOrderCommand(1, 25.50, products)

		// Act
		_, err := useCase.Execute(command)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "products not found")
	})
}

func TestGetOrderUseCase_Execute(t *testing.T) {
	t.Run("deve buscar pedido com sucesso", func(t *testing.T) {
		// Arrange
		orderRepo := mocks.NewMockOrderRepository(t)

		expectedOrder := &entities.OrderEntity{
			ID:          1,
			CustomerId:  1,
			TotalAmount: 51.00,
			Products: []*entities.OrderProductEntity{
				{ID: 1, OrderId: 1, ProductId: 101, Price: 25.50, Quantity: 2},
			},
			Status: []*entities.OrderStatusEntity{
				{ID: 1, OrderId: 1, CurrentStatus: 1},
			},
		}

		orderRepo.On("GetOrder", uint(1)).Return(expectedOrder, nil)

		useCase := getorder.NewGetOrderUseCaseImpl(orderRepo)
		command := commands.NewGetOrderCommand(1)

		// Act
		result, err := useCase.Execute(command)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, float32(51.00), result.TotalAmount)
		orderRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro quando pedido não existe", func(t *testing.T) {
		// Arrange
		orderRepo := mocks.NewMockOrderRepository(t)
		orderRepo.On("GetOrder", uint(999)).Return(nil, assert.AnError)

		useCase := getorder.NewGetOrderUseCaseImpl(orderRepo)
		command := commands.NewGetOrderCommand(999)

		// Act
		_, err := useCase.Execute(command)

		// Assert
		require.Error(t, err)
	})
}

func TestGetOrdersUseCase_Execute(t *testing.T) {
	t.Run("deve listar pedidos com sucesso", func(t *testing.T) {
		// Arrange
		orderRepo := mocks.NewMockOrderRepository(t)

		expectedOrders := []*entities.OrderEntity{
			{ID: 1, CustomerId: 1, TotalAmount: 51.00},
			{ID: 2, CustomerId: 2, TotalAmount: 30.00},
		}

		orderRepo.On("GetOrders").Return(expectedOrders, nil)

		useCase := getorders.NewGetOrdersUseCaseImpl(orderRepo)
		command := commands.NewGetOrdersCommand()

		// Act
		result, err := useCase.Execute(command)

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, uint(1), result[0].ID)
		assert.Equal(t, uint(2), result[1].ID)
		orderRepo.AssertExpectations(t)
	})
}

func TestGetOrderStatusUseCase_Execute(t *testing.T) {
	t.Run("deve buscar status do pedido", func(t *testing.T) {
		// Arrange
		statusRepo := mocks.NewMockOrderStatusRepository(t)

		expectedStatus := &entities.OrderStatusEntity{
			ID:            1,
			OrderId:       1,
			CurrentStatus: 2,
		}

		statusRepo.On("GetOrderStatus", uint(1)).Return(expectedStatus, nil)

		useCase := getorderstatus.NewGetOrderStatusUseCaseImpl(statusRepo)
		command := commands.NewGetOrderStatusCommand(1)

		// Act
		result, err := useCase.Execute(command)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, uint(1), result.OrderId)
		assert.Equal(t, uint(2), result.CurrentStatus)
		statusRepo.AssertExpectations(t)
	})
}

func TestUpdateOrderStatusUseCase_Execute(t *testing.T) {
	t.Run("deve atualizar status do pedido", func(t *testing.T) {
		// Arrange
		statusRepo := mocks.NewMockOrderStatusRepository(t)

		statusRepo.On("AddOrderStatus", mock.MatchedBy(func(status *entities.OrderStatusEntity) bool {
			return status.OrderId == 1 && status.CurrentStatus == 2
		})).Return(nil)

		useCase := updateorderstatus.NewUpdateOrderStatusUseCaseImpl(statusRepo)
		command := commands.NewUpdateOrderStatusCommand(1, 2)

		// Act
		err := useCase.Execute(command)

		// Assert
		require.NoError(t, err)
		statusRepo.AssertExpectations(t)
	})
}
