package updateorderstatus_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
	updateorderstatus "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/updateOrderStatus"
	mockRepositories "github.com/viniciuscluna/tc-fiap-50/mocks/order/domain/repositories"
)

type UpdateOrderStatusUseCaseTestSuite struct {
	suite.Suite
	mockOrderStatusRepository *mockRepositories.MockOrderStatusRepository
	useCase                   updateorderstatus.UpdateOrderStatusUseCase
}

func (suite *UpdateOrderStatusUseCaseTestSuite) SetupTest() {
	suite.mockOrderStatusRepository = mockRepositories.NewMockOrderStatusRepository(suite.T())
	suite.useCase = updateorderstatus.NewUpdateOrderStatusUseCaseImpl(suite.mockOrderStatusRepository)
}

func TestUpdateOrderStatusUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateOrderStatusUseCaseTestSuite))
}

// Feature: Update Order Status Use Case
// Scenario: Update order status successfully

func (suite *UpdateOrderStatusUseCaseTestSuite) Test_UpdateOrderStatus_WithValidCommand_ShouldAddNewStatus() {
	// GIVEN a valid update order status command
	orderId := uint(123)
	newStatus := uint(2)
	command := commands.NewUpdateOrderStatusCommand(orderId, newStatus)

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.MatchedBy(func(status *entities.OrderStatusEntity) bool {
			return status.OrderId == orderId && status.CurrentStatus == newStatus
		})).
		Return(nil).
		Once()

	// WHEN the order status is updated
	err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND the new status should have been added
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *UpdateOrderStatusUseCaseTestSuite) Test_UpdateOrderStatus_ToEmPreparacao_ShouldSetStatus2() {
	// GIVEN an order transitioning to "Em preparação" (status 2)
	orderId := uint(100)
	command := commands.NewUpdateOrderStatusCommand(orderId, 2)

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.MatchedBy(func(status *entities.OrderStatusEntity) bool {
			return status.OrderId == 100 && status.CurrentStatus == 2
		})).
		Return(nil).
		Once()

	// WHEN the status is updated
	err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *UpdateOrderStatusUseCaseTestSuite) Test_UpdateOrderStatus_ToPronto_ShouldSetStatus3() {
	// GIVEN an order transitioning to "Pronto" (status 3)
	orderId := uint(200)
	command := commands.NewUpdateOrderStatusCommand(orderId, 3)

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.MatchedBy(func(status *entities.OrderStatusEntity) bool {
			return status.OrderId == 200 && status.CurrentStatus == 3
		})).
		Return(nil).
		Once()

	// WHEN the status is updated
	err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *UpdateOrderStatusUseCaseTestSuite) Test_UpdateOrderStatus_ToFinalizado_ShouldSetStatus4() {
	// GIVEN an order transitioning to "Finalizado" (status 4)
	orderId := uint(300)
	command := commands.NewUpdateOrderStatusCommand(orderId, 4)

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.MatchedBy(func(status *entities.OrderStatusEntity) bool {
			return status.OrderId == 300 && status.CurrentStatus == 4
		})).
		Return(nil).
		Once()

	// WHEN the status is updated
	err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *UpdateOrderStatusUseCaseTestSuite) Test_UpdateOrderStatus_WithRepositoryError_ShouldReturnError() {
	// GIVEN a valid command
	orderId := uint(400)
	command := commands.NewUpdateOrderStatusCommand(orderId, 2)

	// AND the repository encounters an error
	expectedError := errors.New("database connection error")

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.Anything).
		Return(expectedError).
		Once()

	// WHEN attempting to update the status
	err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	// AND the error should match the repository error
	assert.Equal(suite.T(), expectedError, err)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *UpdateOrderStatusUseCaseTestSuite) Test_UpdateOrderStatus_WithInvalidOrderId_ShouldReturnError() {
	// GIVEN a command with non-existent order ID
	command := commands.NewUpdateOrderStatusCommand(9999, 2)

	expectedError := errors.New("foreign key constraint failed")

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.Anything).
		Return(expectedError).
		Once()

	// WHEN attempting to update the status
	err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}

func (suite *UpdateOrderStatusUseCaseTestSuite) Test_UpdateOrderStatus_ShouldCreateNewStatusRecord() {
	// GIVEN an order with existing status history
	orderId := uint(500)
	command := commands.NewUpdateOrderStatusCommand(orderId, 3)

	suite.mockOrderStatusRepository.EXPECT().
		AddOrderStatus(mock.MatchedBy(func(status *entities.OrderStatusEntity) bool {
			return status.OrderId == 500 && status.CurrentStatus == 3
		})).
		Return(nil).
		Once()

	// WHEN the status is updated
	err := suite.useCase.Execute(command)

	// THEN a new status record should be created (not updated)
	assert.NoError(suite.T(), err)
	suite.mockOrderStatusRepository.AssertExpectations(suite.T())
}
