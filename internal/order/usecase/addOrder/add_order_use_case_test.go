package addorder_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"
	addorder "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/addOrder"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
	mockClients "github.com/viniciuscluna/tc-fiap-50/mocks/infrastructure/clients"
	mockRepositories "github.com/viniciuscluna/tc-fiap-50/mocks/order/domain/repositories"
)

type AddOrderUseCaseTestSuite struct {
	suite.Suite
	mockOrderRepository        *mockRepositories.MockOrderRepository
	mockOrderProductRepository *mockRepositories.MockOrderProductRepository
	mockOrderStatusRepository  *mockRepositories.MockOrderStatusRepository
	mockCustomerClient         *mockClients.MockCustomerClient
	mockProductClient          *mockClients.MockProductClient
	useCase                    addorder.AddOrderUseCase
}

func (suite *AddOrderUseCaseTestSuite) SetupTest() {
	suite.mockOrderRepository = mockRepositories.NewMockOrderRepository(suite.T())
	suite.mockOrderProductRepository = mockRepositories.NewMockOrderProductRepository(suite.T())
	suite.mockOrderStatusRepository = mockRepositories.NewMockOrderStatusRepository(suite.T())
	suite.mockCustomerClient = mockClients.NewMockCustomerClient(suite.T())
	suite.mockProductClient = mockClients.NewMockProductClient(suite.T())

	suite.useCase = addorder.NewAddOrderUseCaseImpl(
		suite.mockOrderRepository,
		suite.mockOrderProductRepository,
		suite.mockOrderStatusRepository,
		suite.mockCustomerClient,
		suite.mockProductClient,
	)
}

func TestAddOrderUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AddOrderUseCaseTestSuite))
}

// Feature: Add Order Use Case
// Scenario: Create a new order with products successfully

func (suite *AddOrderUseCaseTestSuite) Test_AddOrder_WithValidCommand_ShouldCreateOrderSuccessfully() {
	// GIVEN a valid add order command with products
	products := []*dto.AddOrderProductDto{
		{ProductId: 1, Quantity: 2, Price: 25.50},
		{ProductId: 2, Quantity: 1, Price: 50.00},
	}
	command := commands.NewAddOrderCommand(1, 101.00, products)

	createdOrder := &entities.OrderEntity{
		ID:          123,
		CustomerId:  1,
		TotalAmount: 101.00,
	}

	suite.mockOrderRepository.EXPECT().
		AddOrder(mock.MatchedBy(func(order *entities.OrderEntity) bool {
			return order.CustomerId == 1 && order.TotalAmount == 101.00
		})).
		Return(createdOrder, nil).
		Once()

	suite.mockOrderProductRepository.EXPECT().
		AddOrderProduct(mock.MatchedBy(func(product *entities.OrderProductEntity) bool {
			return product.OrderId == 123 && product.ProductId == 1
		})).
		Return(nil).
		Once()

	suite.mockOrderProductRepository.EXPECT().
		AddOrderProduct(mock.MatchedBy(func(product *entities.OrderProductEntity) bool {
			return product.OrderId == 123 && product.ProductId == 2
		})).
		Return(nil).
		Once()

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.MatchedBy(func(status *entities.OrderStatusEntity) bool {
			return status.OrderId == 123 && status.CurrentStatus == 1
		})).
		Return(nil).
		Once()

	// WHEN the order creation is executed
	orderId, err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND the order ID should be returned
	assert.Equal(suite.T(), "123", orderId)
	// AND all repository methods should have been called
	suite.mockOrderRepository.AssertExpectations(suite.T())
	suite.mockOrderProductRepository.AssertExpectations(suite.T())
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *AddOrderUseCaseTestSuite) Test_AddOrder_WithMultipleProducts_ShouldAddAllProducts() {
	// GIVEN an order with three different products
	products := []*dto.AddOrderProductDto{
		{ProductId: 10, Quantity: 1, Price: 10.00},
		{ProductId: 20, Quantity: 2, Price: 20.00},
		{ProductId: 30, Quantity: 3, Price: 30.00},
	}
	command := commands.NewAddOrderCommand(5, 140.00, products)

	createdOrder := &entities.OrderEntity{ID: 456, CustomerId: 5, TotalAmount: 140.00}

	suite.mockOrderRepository.EXPECT().
		AddOrder(mock.Anything).
		Return(createdOrder, nil).
		Once()

	suite.mockOrderProductRepository.EXPECT().
		AddOrderProduct(mock.Anything).
		Return(nil).
		Times(3)

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.Anything).
		Return(nil).
		Once()

	// WHEN the order is created
	orderId, err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "456", orderId)
	// AND all three products should have been added
	suite.mockOrderProductRepository.AssertNumberOfCalls(suite.T(), "AddOrderProduct", 3)
}

func (suite *AddOrderUseCaseTestSuite) Test_AddOrder_ShouldSetInitialStatusToRecebido() {
	// GIVEN a valid order command
	products := []*dto.AddOrderProductDto{
		{ProductId: 1, Quantity: 1, Price: 50.00},
	}
	command := commands.NewAddOrderCommand(1, 50.00, products)

	createdOrder := &entities.OrderEntity{ID: 789, CustomerId: 1, TotalAmount: 50.00}

	suite.mockOrderRepository.EXPECT().
		AddOrder(mock.Anything).
		Return(createdOrder, nil).
		Once()

	suite.mockOrderProductRepository.EXPECT().
		AddOrderProduct(mock.Anything).
		Return(nil).
		Once()

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.MatchedBy(func(status *entities.OrderStatusEntity) bool {
			return status.OrderId == 789 && status.CurrentStatus == 1
		})).
		Return(nil).
		Once()

	// WHEN the order is created
	orderId, err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "789", orderId)
	// AND the initial status should be 1 (Recebido)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *AddOrderUseCaseTestSuite) Test_AddOrder_WithOrderRepositoryError_ShouldReturnError() {
	// GIVEN a valid order command
	products := []*dto.AddOrderProductDto{
		{ProductId: 1, Quantity: 1, Price: 50.00},
	}
	command := commands.NewAddOrderCommand(1, 50.00, products)

	expectedError := errors.New("database connection error")

	suite.mockOrderRepository.EXPECT().
		AddOrder(mock.Anything).
		Return(nil, expectedError).
		Once()

	// WHEN the order creation is attempted
	orderId, err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	// AND the error should match the repository error
	assert.Equal(suite.T(), expectedError, err)
	// AND no order ID should be returned
	assert.Empty(suite.T(), orderId)
	// AND products should not have been added
	suite.mockOrderProductRepository.AssertNotCalled(suite.T(), "AddOrderProduct")
	// AND status should not have been added
	suite.mockOrderStatusRepository.AssertNotCalled(suite.T(), "AddOrderStatus")
}

func (suite *AddOrderUseCaseTestSuite) Test_AddOrder_WithProductRepositoryError_ShouldReturnError() {
	// GIVEN a valid order command
	products := []*dto.AddOrderProductDto{
		{ProductId: 1, Quantity: 1, Price: 50.00},
	}
	command := commands.NewAddOrderCommand(1, 50.00, products)

	createdOrder := &entities.OrderEntity{ID: 100, CustomerId: 1, TotalAmount: 50.00}
	expectedError := errors.New("product insert error")

	suite.mockOrderRepository.EXPECT().
		AddOrder(mock.Anything).
		Return(createdOrder, nil).
		Once()

	suite.mockOrderProductRepository.EXPECT().
		AddOrderProduct(mock.Anything).
		Return(expectedError).
		Once()

	// WHEN the order creation is attempted
	orderId, err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	// AND no order ID should be returned
	assert.Empty(suite.T(), orderId)
	// AND status should not have been added
	suite.mockOrderStatusRepository.AssertNotCalled(suite.T(), "AddOrderStatus")
}

func (suite *AddOrderUseCaseTestSuite) Test_AddOrder_WithStatusRepositoryError_ShouldReturnError() {
	// GIVEN a valid order command
	products := []*dto.AddOrderProductDto{
		{ProductId: 1, Quantity: 1, Price: 50.00},
	}
	command := commands.NewAddOrderCommand(1, 50.00, products)

	createdOrder := &entities.OrderEntity{ID: 200, CustomerId: 1, TotalAmount: 50.00}
	expectedError := errors.New("status insert error")

	suite.mockOrderRepository.EXPECT().
		AddOrder(mock.Anything).
		Return(createdOrder, nil).
		Once()

	suite.mockOrderProductRepository.EXPECT().
		AddOrderProduct(mock.Anything).
		Return(nil).
		Once()

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.Anything).
		Return(expectedError).
		Once()

	// WHEN the order creation is attempted
	orderId, err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	// AND no order ID should be returned
	assert.Empty(suite.T(), orderId)
}

func (suite *AddOrderUseCaseTestSuite) Test_AddOrder_WithNoProducts_ShouldCreateOrderWithoutProducts() {
	// GIVEN an order command without products
	products := []*dto.AddOrderProductDto{}
	command := commands.NewAddOrderCommand(1, 0, products)

	createdOrder := &entities.OrderEntity{ID: 300, CustomerId: 1, TotalAmount: 0}

	suite.mockOrderRepository.EXPECT().
		AddOrder(mock.Anything).
		Return(createdOrder, nil).
		Once()

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.Anything).
		Return(nil).
		Once()

	// WHEN the order is created
	orderId, err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "300", orderId)
	// AND no products should have been added
	suite.mockOrderProductRepository.AssertNotCalled(suite.T(), "AddOrderProduct")
	// AND the status should still be created
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}
