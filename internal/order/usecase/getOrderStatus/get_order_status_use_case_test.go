package getorderstatus_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
	getorderstatus "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrderStatus"
	mockRepositories "github.com/viniciuscluna/tc-fiap-50/mocks/order/domain/repositories"
)

type GetOrderStatusUseCaseTestSuite struct {
	suite.Suite
	mockOrderStatusRepository *mockRepositories.MockOrderStatusRepository
	useCase                   getorderstatus.GetOrderStatusUseCase
}

func (suite *GetOrderStatusUseCaseTestSuite) SetupTest() {
	suite.mockOrderStatusRepository = mockRepositories.NewMockOrderStatusRepository(suite.T())
	suite.useCase = getorderstatus.NewGetOrderStatusUseCaseImpl(suite.mockOrderStatusRepository)
}

func TestGetOrderStatusUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrderStatusUseCaseTestSuite))
}

// Feature: Get Order Status Use Case
// Scenario: Retrieve current status for an order

func (suite *GetOrderStatusUseCaseTestSuite) Test_GetOrderStatus_WithValidOrderId_ShouldReturnLatestStatus() {
	// GIVEN a valid order ID with status history
	orderId := uint(123)
	command := commands.NewGetOrderStatusCommand(orderId)

	expectedStatus := &entities.OrderStatusEntity{
		ID:            10,
		OrderId:       123,
		CurrentStatus: 2,
		CreatedAt:     time.Now(),
	}

	suite.mockOrderStatusRepository.EXPECT().
		GetOrderStatus(orderId).
		Return(expectedStatus, nil).
		Once()

	// WHEN the order status is retrieved
	result, err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// AND the latest status should be returned
	assert.Equal(suite.T(), expectedStatus.ID, result.ID)
	assert.Equal(suite.T(), expectedStatus.OrderId, result.OrderId)
	assert.Equal(suite.T(), expectedStatus.CurrentStatus, result.CurrentStatus)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *GetOrderStatusUseCaseTestSuite) Test_GetOrderStatus_WithNonExistentOrderId_ShouldReturnError() {
	// GIVEN a non-existent order ID
	orderId := uint(9999)
	command := commands.NewGetOrderStatusCommand(orderId)

	expectedError := errors.New("record not found")

	suite.mockOrderStatusRepository.EXPECT().
		GetOrderStatus(orderId).
		Return(nil, expectedError).
		Once()

	// WHEN attempting to retrieve the status
	result, err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	// AND no status should be returned
	assert.Nil(suite.T(), result)
	// AND the error should match the repository error
	assert.Equal(suite.T(), expectedError, err)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *GetOrderStatusUseCaseTestSuite) Test_GetOrderStatus_WithRepositoryError_ShouldReturnError() {
	// GIVEN a valid order ID
	orderId := uint(100)
	command := commands.NewGetOrderStatusCommand(orderId)

	// AND the repository encounters an error
	expectedError := errors.New("database connection error")

	suite.mockOrderStatusRepository.EXPECT().
		GetOrderStatus(orderId).
		Return(nil, expectedError).
		Once()

	// WHEN attempting to retrieve the status
	result, err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *GetOrderStatusUseCaseTestSuite) Test_GetOrderStatus_WithRecebidoStatus_ShouldReturnStatus1() {
	// GIVEN an order with "Recebido" status
	orderId := uint(200)
	command := commands.NewGetOrderStatusCommand(orderId)

	expectedStatus := &entities.OrderStatusEntity{
		ID:            1,
		OrderId:       200,
		CurrentStatus: 1,
		CreatedAt:     time.Now(),
	}

	suite.mockOrderStatusRepository.EXPECT().
		GetOrderStatus(orderId).
		Return(expectedStatus, nil).
		Once()

	// WHEN the status is retrieved
	result, err := suite.useCase.Execute(command)

	// THEN the status should be 1 (Recebido)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(1), result.CurrentStatus)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *GetOrderStatusUseCaseTestSuite) Test_GetOrderStatus_WithFinalizadoStatus_ShouldReturnStatus4() {
	// GIVEN an order with "Finalizado" status
	orderId := uint(300)
	command := commands.NewGetOrderStatusCommand(orderId)

	expectedStatus := &entities.OrderStatusEntity{
		ID:            5,
		OrderId:       300,
		CurrentStatus: 4,
		CreatedAt:     time.Now(),
	}

	suite.mockOrderStatusRepository.EXPECT().
		GetOrderStatus(orderId).
		Return(expectedStatus, nil).
		Once()

	// WHEN the status is retrieved
	result, err := suite.useCase.Execute(command)

	// THEN the status should be 4 (Finalizado)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(4), result.CurrentStatus)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}
