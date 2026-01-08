package getorder_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
	getorder "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrder"
	mockRepositories "github.com/viniciuscluna/tc-fiap-50/mocks/order/domain/repositories"
)

type GetOrderUseCaseTestSuite struct {
	suite.Suite
	mockOrderRepository *mockRepositories.MockOrderRepository
	useCase             getorder.GetOrderUseCase
}

func (suite *GetOrderUseCaseTestSuite) SetupTest() {
	suite.mockOrderRepository = mockRepositories.NewMockOrderRepository(suite.T())
	suite.useCase = getorder.NewGetOrderUseCaseImpl(suite.mockOrderRepository)
}

func TestGetOrderUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrderUseCaseTestSuite))
}

// Feature: Get Order Use Case
// Scenario: Retrieve an order by ID successfully

func (suite *GetOrderUseCaseTestSuite) Test_GetOrder_WithValidOrderId_ShouldReturnOrder() {
	// GIVEN a valid order ID and an existing order
	orderId := uint(123)
	command := commands.NewGetOrderCommand(orderId)

	expectedOrder := &entities.OrderEntity{
		ID:          123,
		CustomerId:  1,
		TotalAmount: 100.50,
		CreatedAt:   time.Now(),
		Products: []*entities.OrderProductEntity{
			{ID: 1, ProductId: 10, Quantity: 2, Price: 50.25},
		},
		Status: []*entities.OrderStatusEntity{
			{ID: 1, CurrentStatus: 2, OrderId: 123},
		},
	}

	suite.mockOrderRepository.EXPECT().
		GetOrder(orderId).
		Return(expectedOrder, nil).
		Once()

	// WHEN the order is retrieved
	result, err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// AND the order data should match
	assert.Equal(suite.T(), expectedOrder.ID, result.ID)
	assert.Equal(suite.T(), expectedOrder.CustomerId, result.CustomerId)
	assert.Equal(suite.T(), expectedOrder.TotalAmount, result.TotalAmount)
	// AND related data should be included
	assert.NotNil(suite.T(), result.Products)
	assert.Len(suite.T(), result.Products, 1)
	assert.NotNil(suite.T(), result.Status)
	assert.Len(suite.T(), result.Status, 1)
	suite.mockOrderRepository.AssertExpectations(suite.T())
}

func (suite *GetOrderUseCaseTestSuite) Test_GetOrder_WithNonExistentOrderId_ShouldReturnError() {
	// GIVEN a non-existent order ID
	orderId := uint(9999)
	command := commands.NewGetOrderCommand(orderId)

	expectedError := errors.New("record not found")

	suite.mockOrderRepository.EXPECT().
		GetOrder(orderId).
		Return(nil, expectedError).
		Once()

	// WHEN attempting to retrieve the order
	result, err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	// AND no order should be returned
	assert.Nil(suite.T(), result)
	// AND the error should match the repository error
	assert.Equal(suite.T(), expectedError, err)
	suite.mockOrderRepository.AssertExpectations(suite.T())
}

func (suite *GetOrderUseCaseTestSuite) Test_GetOrder_WithRepositoryError_ShouldReturnError() {
	// GIVEN a valid order ID
	orderId := uint(100)
	command := commands.NewGetOrderCommand(orderId)

	// AND a repository error occurs
	expectedError := errors.New("database connection error")

	suite.mockOrderRepository.EXPECT().
		GetOrder(orderId).
		Return(nil, expectedError).
		Once()

	// WHEN attempting to retrieve the order
	result, err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockOrderRepository.AssertExpectations(suite.T())
}

func (suite *GetOrderUseCaseTestSuite) Test_GetOrder_WithOrderWithoutProducts_ShouldReturnOrderSuccessfully() {
	// GIVEN an order without products
	orderId := uint(200)
	command := commands.NewGetOrderCommand(orderId)

	expectedOrder := &entities.OrderEntity{
		ID:          200,
		CustomerId:  5,
		TotalAmount: 0,
		Products:    []*entities.OrderProductEntity{},
		Status: []*entities.OrderStatusEntity{
			{ID: 1, CurrentStatus: 1, OrderId: 200},
		},
	}

	suite.mockOrderRepository.EXPECT().
		GetOrder(orderId).
		Return(expectedOrder, nil).
		Once()

	// WHEN the order is retrieved
	result, err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// AND the order should have empty products list
	assert.Empty(suite.T(), result.Products)
	// AND the status should still be present
	assert.NotEmpty(suite.T(), result.Status)
	suite.mockOrderRepository.AssertExpectations(suite.T())
}
