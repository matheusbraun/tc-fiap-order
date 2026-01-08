package controller_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/controller"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"
	mockPresenter "github.com/viniciuscluna/tc-fiap-50/mocks/order/presenter"
	mockAddOrder "github.com/viniciuscluna/tc-fiap-50/mocks/order/usecase/addOrder"
	mockGetOrder "github.com/viniciuscluna/tc-fiap-50/mocks/order/usecase/getOrder"
	mockGetOrderStatus "github.com/viniciuscluna/tc-fiap-50/mocks/order/usecase/getOrderStatus"
	mockGetOrders "github.com/viniciuscluna/tc-fiap-50/mocks/order/usecase/getOrders"
	mockUpdateOrderStatus "github.com/viniciuscluna/tc-fiap-50/mocks/order/usecase/updateOrderStatus"
)

type OrderControllerTestSuite struct {
	suite.Suite
	mockPresenter                *mockPresenter.MockOrderPresenter
	mockAddOrderUseCase          *mockAddOrder.MockAddOrderUseCase
	mockGetOrderUseCase          *mockGetOrder.MockGetOrderUseCase
	mockGetOrdersUseCase         *mockGetOrders.MockGetOrdersUseCase
	mockGetOrderStatusUseCase    *mockGetOrderStatus.MockGetOrderStatusUseCase
	mockUpdateOrderStatusUseCase *mockUpdateOrderStatus.MockUpdateOrderStatusUseCase
	controller                   controller.OrderController
}

func (suite *OrderControllerTestSuite) SetupTest() {
	suite.mockPresenter = mockPresenter.NewMockOrderPresenter(suite.T())
	suite.mockAddOrderUseCase = mockAddOrder.NewMockAddOrderUseCase(suite.T())
	suite.mockGetOrderUseCase = mockGetOrder.NewMockGetOrderUseCase(suite.T())
	suite.mockGetOrdersUseCase = mockGetOrders.NewMockGetOrdersUseCase(suite.T())
	suite.mockGetOrderStatusUseCase = mockGetOrderStatus.NewMockGetOrderStatusUseCase(suite.T())
	suite.mockUpdateOrderStatusUseCase = mockUpdateOrderStatus.NewMockUpdateOrderStatusUseCase(suite.T())

	suite.controller = controller.NewOrderControllerImpl(
		suite.mockPresenter,
		suite.mockAddOrderUseCase,
		suite.mockGetOrderUseCase,
		suite.mockGetOrdersUseCase,
		suite.mockGetOrderStatusUseCase,
		suite.mockUpdateOrderStatusUseCase,
	)
}

func TestOrderControllerTestSuite(t *testing.T) {
	suite.Run(t, new(OrderControllerTestSuite))
}

// Feature: Order Controller - Add Order
// Scenario: Create a new order successfully

func (suite *OrderControllerTestSuite) Test_Add_WithValidDto_ShouldCreateOrderSuccessfully() {
	// GIVEN a valid add order DTO
	customerId := uint(1)
	addOrderDto := &dto.AddOrderDto{
		CustomerId:  &customerId,
		TotalAmount: 100.00,
		Products: []*dto.AddOrderProductDto{
			{ProductId: 10, Quantity: 2, Price: 50.00},
		},
	}

	expectedOrderId := "123"

	suite.mockAddOrderUseCase.EXPECT().
		Execute(mock.Anything).
		Return(expectedOrderId, nil).
		Once()

	// WHEN the order is added
	orderId, err := suite.controller.Add(addOrderDto)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND the order ID should be returned
	assert.Equal(suite.T(), expectedOrderId, orderId)
	suite.mockAddOrderUseCase.AssertExpectations(suite.T())
}

func (suite *OrderControllerTestSuite) Test_Add_WithUseCaseError_ShouldReturnError() {
	// GIVEN a valid add order DTO
	customerId := uint(1)
	addOrderDto := &dto.AddOrderDto{
		CustomerId:  &customerId,
		TotalAmount: 50.00,
		Products:    []*dto.AddOrderProductDto{},
	}

	// AND the use case returns an error
	expectedError := errors.New("database error")

	suite.mockAddOrderUseCase.EXPECT().
		Execute(mock.Anything).
		Return("", expectedError).
		Once()

	// WHEN attempting to add the order
	orderId, err := suite.controller.Add(addOrderDto)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	// AND no order ID should be returned
	assert.Empty(suite.T(), orderId)
	suite.mockAddOrderUseCase.AssertExpectations(suite.T())
}

// Feature: Order Controller - Get Order
// Scenario: Retrieve and present order data

func (suite *OrderControllerTestSuite) Test_GetOrder_WithValidId_ShouldReturnPresentedOrder() {
	// GIVEN a valid order ID
	orderId := uint(123)

	orderEntity := &entities.OrderEntity{
		ID:          123,
		CustomerId:  1,
		TotalAmount: 150.00,
		CreatedAt:   time.Now(),
	}

	expectedDto := &dto.GetOrderResponseDto{
		ID:          123,
		CustomerId:  1,
		TotalAmount: 150.00,
	}

	suite.mockGetOrderUseCase.EXPECT().
		Execute(mock.Anything).
		Return(orderEntity, nil).
		Once()

	suite.mockPresenter.EXPECT().
		Present(orderEntity).
		Return(expectedDto).
		Once()

	// WHEN the order is retrieved
	result, err := suite.controller.GetOrder(orderId)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// AND the order data should be correctly presented
	assert.Equal(suite.T(), expectedDto.ID, result.ID)
	assert.Equal(suite.T(), expectedDto.CustomerId, result.CustomerId)
	assert.Equal(suite.T(), expectedDto.TotalAmount, result.TotalAmount)
	suite.mockGetOrderUseCase.AssertExpectations(suite.T())
	suite.mockPresenter.AssertExpectations(suite.T())
}

func (suite *OrderControllerTestSuite) Test_GetOrder_WithNonExistentId_ShouldReturnError() {
	// GIVEN a non-existent order ID
	orderId := uint(9999)
	expectedError := errors.New("order not found")

	suite.mockGetOrderUseCase.EXPECT().
		Execute(mock.Anything).
		Return(nil, expectedError).
		Once()

	// WHEN attempting to retrieve the order
	result, err := suite.controller.GetOrder(orderId)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), expectedError, err)
	// AND the presenter should not be called
	suite.mockPresenter.AssertNotCalled(suite.T(), "Present")
	suite.mockGetOrderUseCase.AssertExpectations(suite.T())
}

// Feature: Order Controller - Get Orders
// Scenario: List all active orders

func (suite *OrderControllerTestSuite) Test_GetOrders_ShouldReturnPresentedOrdersList() {
	// GIVEN multiple active orders
	orders := []*entities.OrderEntity{
		{ID: 1, CustomerId: 1, TotalAmount: 50.00},
		{ID: 2, CustomerId: 2, TotalAmount: 75.00},
	}

	expectedDto := &dto.GetOrdersResponseDto{
		Orders: []*dto.GetOrderResponseDto{
			{ID: 1, CustomerId: 1, TotalAmount: 50.00},
			{ID: 2, CustomerId: 2, TotalAmount: 75.00},
		},
	}

	suite.mockGetOrdersUseCase.EXPECT().
		Execute(mock.Anything).
		Return(orders, nil).
		Once()

	suite.mockPresenter.EXPECT().
		PresentOrders(orders).
		Return(expectedDto).
		Once()

	// WHEN all orders are retrieved
	result, err := suite.controller.GetOrders()

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// AND the orders should be correctly presented
	assert.Len(suite.T(), result.Orders, 2)
	suite.mockGetOrdersUseCase.AssertExpectations(suite.T())
	suite.mockPresenter.AssertExpectations(suite.T())
}

func (suite *OrderControllerTestSuite) Test_GetOrders_WithUseCaseError_ShouldReturnError() {
	// GIVEN the use case returns an error
	expectedError := errors.New("database connection error")

	suite.mockGetOrdersUseCase.EXPECT().
		Execute(mock.Anything).
		Return(nil, expectedError).
		Once()

	// WHEN attempting to retrieve orders
	result, err := suite.controller.GetOrders()

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), expectedError, err)
	// AND the presenter should not be called
	suite.mockPresenter.AssertNotCalled(suite.T(), "PresentOrders")
	suite.mockGetOrdersUseCase.AssertExpectations(suite.T())
}

// Feature: Order Controller - Get Order Status
// Scenario: Retrieve order status

func (suite *OrderControllerTestSuite) Test_GetOrderStatus_WithValidId_ShouldReturnPresentedStatus() {
	// GIVEN a valid order ID
	orderId := uint(100)

	statusEntity := &entities.OrderStatusEntity{
		ID:            1,
		OrderId:       100,
		CurrentStatus: 2,
		CreatedAt:     time.Now(),
	}

	expectedDto := &dto.GetOrderStatusResponseDto{
		ID:            1,
		OrderId:       100,
		CurrentStatus: 2,
	}

	suite.mockGetOrderStatusUseCase.EXPECT().
		Execute(mock.Anything).
		Return(statusEntity, nil).
		Once()

	suite.mockPresenter.EXPECT().
		PresentStatus(statusEntity).
		Return(expectedDto).
		Once()

	// WHEN the order status is retrieved
	result, err := suite.controller.GetOrderStatus(orderId)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// AND the status should be correctly presented
	assert.Equal(suite.T(), expectedDto.CurrentStatus, result.CurrentStatus)
	suite.mockGetOrderStatusUseCase.AssertExpectations(suite.T())
	suite.mockPresenter.AssertExpectations(suite.T())
}

func (suite *OrderControllerTestSuite) Test_GetOrderStatus_WithError_ShouldReturnError() {
	// GIVEN an order ID
	orderId := uint(200)
	expectedError := errors.New("status not found")

	suite.mockGetOrderStatusUseCase.EXPECT().
		Execute(mock.Anything).
		Return(nil, expectedError).
		Once()

	// WHEN attempting to retrieve the status
	result, err := suite.controller.GetOrderStatus(orderId)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockGetOrderStatusUseCase.AssertExpectations(suite.T())
}

// Feature: Order Controller - Update Order Status
// Scenario: Update order status successfully

func (suite *OrderControllerTestSuite) Test_UpdateOrderStatus_WithValidRequest_ShouldUpdateSuccessfully() {
	// GIVEN a valid order ID and status update request
	orderId := uint(123)
	updateRequest := &dto.UpdateOrderStatusRequestDto{
		Status: 3,
	}

	suite.mockUpdateOrderStatusUseCase.EXPECT().
		Execute(mock.Anything).
		Return(nil).
		Once()

	// WHEN the order status is updated
	err := suite.controller.UpdateOrderStatus(orderId, updateRequest)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	suite.mockUpdateOrderStatusUseCase.AssertExpectations(suite.T())
}

func (suite *OrderControllerTestSuite) Test_UpdateOrderStatus_WithUseCaseError_ShouldReturnError() {
	// GIVEN a valid request
	orderId := uint(456)
	updateRequest := &dto.UpdateOrderStatusRequestDto{
		Status: 2,
	}

	// AND the use case returns an error
	expectedError := errors.New("update failed")

	suite.mockUpdateOrderStatusUseCase.EXPECT().
		Execute(mock.Anything).
		Return(expectedError).
		Once()

	// WHEN attempting to update the status
	err := suite.controller.UpdateOrderStatus(orderId, updateRequest)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockUpdateOrderStatusUseCase.AssertExpectations(suite.T())
}
