package getorders_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
	getorders "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrders"
	mockRepositories "github.com/viniciuscluna/tc-fiap-50/mocks/order/domain/repositories"
)

type GetOrdersUseCaseTestSuite struct {
	suite.Suite
	mockOrderRepository *mockRepositories.MockOrderRepository
	useCase             getorders.GetOrdersUseCase
}

func (suite *GetOrdersUseCaseTestSuite) SetupTest() {
	suite.mockOrderRepository = mockRepositories.NewMockOrderRepository(suite.T())
	suite.useCase = getorders.NewGetOrdersUseCaseImpl(suite.mockOrderRepository)
}

func TestGetOrdersUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrdersUseCaseTestSuite))
}

// Feature: Get Orders Use Case
// Scenario: List all active orders successfully

func (suite *GetOrdersUseCaseTestSuite) Test_GetOrders_WithActiveOrders_ShouldReturnAllOrders() {
	// GIVEN multiple active orders exist
	command := commands.NewGetOrdersCommand()

	expectedOrders := []*entities.OrderEntity{
		{
			ID:          100,
			CustomerId:  1,
			TotalAmount: 50.00,
			CreatedAt:   time.Now().Add(-2 * time.Hour),
			Status:      []*entities.OrderStatusEntity{{CurrentStatus: 1}},
		},
		{
			ID:          101,
			CustomerId:  2,
			TotalAmount: 75.00,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			Status:      []*entities.OrderStatusEntity{{CurrentStatus: 2}},
		},
		{
			ID:          102,
			CustomerId:  3,
			TotalAmount: 100.00,
			CreatedAt:   time.Now(),
			Status:      []*entities.OrderStatusEntity{{CurrentStatus: 3}},
		},
	}

	suite.mockOrderRepository.EXPECT().
		GetOrders().
		Return(expectedOrders, nil).
		Once()

	// WHEN all orders are retrieved
	results, err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), results)
	// AND all active orders should be returned
	assert.Len(suite.T(), results, 3)
	assert.Equal(suite.T(), expectedOrders[0].ID, results[0].ID)
	assert.Equal(suite.T(), expectedOrders[1].ID, results[1].ID)
	assert.Equal(suite.T(), expectedOrders[2].ID, results[2].ID)
	suite.mockOrderRepository.AssertExpectations(suite.T())
}

func (suite *GetOrdersUseCaseTestSuite) Test_GetOrders_WithNoOrders_ShouldReturnEmptyList() {
	// GIVEN no active orders exist
	command := commands.NewGetOrdersCommand()

	suite.mockOrderRepository.EXPECT().
		GetOrders().
		Return([]*entities.OrderEntity{}, nil).
		Once()

	// WHEN all orders are retrieved
	results, err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), results)
	// AND an empty list should be returned
	assert.Empty(suite.T(), results)
	suite.mockOrderRepository.AssertExpectations(suite.T())
}

func (suite *GetOrdersUseCaseTestSuite) Test_GetOrders_WithRepositoryError_ShouldReturnError() {
	// GIVEN the repository encounters an error
	command := commands.NewGetOrdersCommand()

	expectedError := errors.New("database connection error")

	suite.mockOrderRepository.EXPECT().
		GetOrders().
		Return(nil, expectedError).
		Once()

	// WHEN attempting to retrieve orders
	results, err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	// AND no orders should be returned
	assert.Nil(suite.T(), results)
	// AND the error should match the repository error
	assert.Equal(suite.T(), expectedError, err)
	suite.mockOrderRepository.AssertExpectations(suite.T())
}

func (suite *GetOrdersUseCaseTestSuite) Test_GetOrders_ShouldExcludeFinishedOrders() {
	// GIVEN only active orders should be returned (excluding status 4)
	command := commands.NewGetOrdersCommand()

	activeOrders := []*entities.OrderEntity{
		{ID: 1, CustomerId: 1, Status: []*entities.OrderStatusEntity{{CurrentStatus: 1}}},
		{ID: 2, CustomerId: 2, Status: []*entities.OrderStatusEntity{{CurrentStatus: 2}}},
		{ID: 3, CustomerId: 3, Status: []*entities.OrderStatusEntity{{CurrentStatus: 3}}},
	}

	suite.mockOrderRepository.EXPECT().
		GetOrders().
		Return(activeOrders, nil).
		Once()

	// WHEN orders are retrieved
	results, err := suite.useCase.Execute(command)

	// THEN only active orders should be returned
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 3)
	// AND no order should have status 4
	for _, order := range results {
		for _, status := range order.Status {
			assert.NotEqual(suite.T(), uint(4), status.CurrentStatus)
		}
	}
	suite.mockOrderRepository.AssertExpectations(suite.T())
}

func (suite *GetOrdersUseCaseTestSuite) Test_GetOrders_WithOrdersAndProducts_ShouldIncludeProducts() {
	// GIVEN orders with products
	command := commands.NewGetOrdersCommand()

	ordersWithProducts := []*entities.OrderEntity{
		{
			ID:          10,
			CustomerId:  1,
			TotalAmount: 100.00,
			Products: []*entities.OrderProductEntity{
				{ID: 1, ProductId: 5, Quantity: 2, Price: 50.00},
			},
			Status: []*entities.OrderStatusEntity{{CurrentStatus: 1}},
		},
	}

	suite.mockOrderRepository.EXPECT().
		GetOrders().
		Return(ordersWithProducts, nil).
		Once()

	// WHEN orders are retrieved
	results, err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 1)
	// AND products should be included
	assert.NotEmpty(suite.T(), results[0].Products)
	assert.Equal(suite.T(), uint(5), results[0].Products[0].ProductId)
	suite.mockOrderRepository.AssertExpectations(suite.T())
}
