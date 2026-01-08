package presenter_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/viniciuscluna/tc-fiap-50/internal/infrastructure/clients"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/presenter"
	mockClients "github.com/viniciuscluna/tc-fiap-50/mocks/infrastructure/clients"
)

type OrderPresenterTestSuite struct {
	suite.Suite
	mockCustomerClient *mockClients.MockCustomerClient
	mockProductClient  *mockClients.MockProductClient
	presenter          presenter.OrderPresenter
}

func (suite *OrderPresenterTestSuite) SetupTest() {
	suite.mockCustomerClient = mockClients.NewMockCustomerClient(suite.T())
	suite.mockProductClient = mockClients.NewMockProductClient(suite.T())
	suite.presenter = presenter.NewOrderPresenterImpl(suite.mockCustomerClient, suite.mockProductClient)
}

func TestOrderPresenterTestSuite(t *testing.T) {
	suite.Run(t, new(OrderPresenterTestSuite))
}

// Feature: Order Presenter - Present Order
// Scenario: Transform order entity to DTO with enriched data

func (suite *OrderPresenterTestSuite) Test_Present_WithValidOrderAndCustomer_ShouldEnrichCustomerData() {
	// GIVEN an order with customer and products
	now := time.Now()
	order := &entities.OrderEntity{
		ID:          123,
		CustomerId:  1,
		TotalAmount: 100.00,
		CreatedAt:   now,
		Products: []*entities.OrderProductEntity{
			{ProductId: 10, Price: 50.00, Quantity: 2},
		},
		Status: []*entities.OrderStatusEntity{
			{ID: 1, CurrentStatus: 1, OrderId: 123, CreatedAt: now},
		},
	}

	customerData := &clients.CustomerDTO{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
		CPF:   12345678901,
	}

	products := []*clients.ProductDTO{
		{ID: 10, Name: "Product A", Price: 50.00, Category: 1},
	}

	suite.mockCustomerClient.EXPECT().
		GetCustomer(mock.Anything, uint(1)).
		Return(customerData, nil).
		Once()

	suite.mockProductClient.EXPECT().
		GetProducts(mock.Anything, []uint{10}).
		Return(products, nil).
		Once()

	// WHEN the order is presented
	result := suite.presenter.Present(order)

	// THEN the DTO should not be nil
	assert.NotNil(suite.T(), result)
	// AND order data should be preserved
	assert.Equal(suite.T(), order.ID, result.ID)
	assert.Equal(suite.T(), order.TotalAmount, result.TotalAmount)
	// AND customer data should be enriched
	assert.NotNil(suite.T(), result.Customer)
	assert.Equal(suite.T(), customerData.Name, result.Customer.Name)
	assert.Equal(suite.T(), customerData.Email, result.Customer.Email)
	// AND products should be enriched
	assert.Len(suite.T(), result.Products, 1)
	assert.Equal(suite.T(), "Product A", result.Products[0].Name)
	suite.mockCustomerClient.AssertExpectations(suite.T())
	suite.mockProductClient.AssertExpectations(suite.T())
}

func (suite *OrderPresenterTestSuite) Test_Present_WithCustomerServiceFailure_ShouldReturnOrderWithoutCustomer() {
	// GIVEN an order with customer ID
	order := &entities.OrderEntity{
		ID:          456,
		CustomerId:  5,
		TotalAmount: 75.00,
		Products:    []*entities.OrderProductEntity{},
		Status:      []*entities.OrderStatusEntity{},
	}

	// AND the customer service fails
	suite.mockCustomerClient.EXPECT().
		GetCustomer(mock.Anything, uint(5)).
		Return(nil, errors.New("service unavailable")).
		Once()

	suite.mockProductClient.EXPECT().
		GetProducts(mock.Anything, mock.Anything).
		Return([]*clients.ProductDTO{}, nil).
		Maybe()

	// WHEN the order is presented
	result := suite.presenter.Present(order)

	// THEN the operation should complete (graceful degradation)
	assert.NotNil(suite.T(), result)
	// AND customer data should be nil
	assert.Nil(suite.T(), result.Customer)
	// AND other order data should still be present
	assert.Equal(suite.T(), order.ID, result.ID)
	assert.Equal(suite.T(), order.CustomerId, result.CustomerId)
	suite.mockCustomerClient.AssertExpectations(suite.T())
}

func (suite *OrderPresenterTestSuite) Test_Present_WithProductServiceFailure_ShouldReturnOrderWithBasicProducts() {
	// GIVEN an order with products
	order := &entities.OrderEntity{
		ID:          789,
		CustomerId:  0,
		TotalAmount: 150.00,
		Products: []*entities.OrderProductEntity{
			{ProductId: 20, Price: 75.00, Quantity: 2},
		},
		Status: []*entities.OrderStatusEntity{},
	}

	// AND the product service fails
	suite.mockProductClient.EXPECT().
		GetProducts(mock.Anything, []uint{20}).
		Return(nil, errors.New("service down")).
		Once()

	// WHEN the order is presented
	result := suite.presenter.Present(order)

	// THEN the operation should complete (graceful degradation)
	assert.NotNil(suite.T(), result)
	// AND products should be returned without enrichment
	assert.Len(suite.T(), result.Products, 1)
	assert.Equal(suite.T(), uint(20), result.Products[0].ProductId)
	assert.Equal(suite.T(), float32(75.00), result.Products[0].Price)
	// AND enriched fields should be empty
	assert.Empty(suite.T(), result.Products[0].Name)
	assert.Empty(suite.T(), result.Products[0].Description)
	suite.mockProductClient.AssertExpectations(suite.T())
}

func (suite *OrderPresenterTestSuite) Test_Present_WithNoCustomerId_ShouldNotFetchCustomer() {
	// GIVEN an order without customer ID
	order := &entities.OrderEntity{
		ID:          100,
		CustomerId:  0,
		TotalAmount: 50.00,
		Products:    []*entities.OrderProductEntity{},
		Status:      []*entities.OrderStatusEntity{},
	}

	suite.mockProductClient.EXPECT().
		GetProducts(mock.Anything, mock.Anything).
		Return([]*clients.ProductDTO{}, nil).
		Maybe()

	// WHEN the order is presented
	result := suite.presenter.Present(order)

	// THEN customer service should not be called
	suite.mockCustomerClient.AssertNotCalled(suite.T(), "GetCustomer")
	// AND customer should be nil in response
	assert.Nil(suite.T(), result.Customer)
	assert.NotNil(suite.T(), result)
}

// Feature: Order Presenter - Present Orders
// Scenario: Transform multiple orders to DTOs

func (suite *OrderPresenterTestSuite) Test_PresentOrders_WithMultipleOrders_ShouldPresentAll() {
	// GIVEN multiple orders
	orders := []*entities.OrderEntity{
		{ID: 1, CustomerId: 1, TotalAmount: 50.00, Products: []*entities.OrderProductEntity{}, Status: []*entities.OrderStatusEntity{}},
		{ID: 2, CustomerId: 2, TotalAmount: 75.00, Products: []*entities.OrderProductEntity{}, Status: []*entities.OrderStatusEntity{}},
	}

	suite.mockCustomerClient.EXPECT().
		GetCustomer(mock.Anything, mock.Anything).
		Return(&clients.CustomerDTO{ID: 1, Name: "Customer"}, nil).
		Times(2)

	suite.mockProductClient.EXPECT().
		GetProducts(mock.Anything, mock.Anything).
		Return([]*clients.ProductDTO{}, nil).
		Times(2)

	// WHEN orders are presented
	result := suite.presenter.PresentOrders(orders)

	// THEN all orders should be in the result
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result.Orders, 2)
	assert.Equal(suite.T(), uint(1), result.Orders[0].ID)
	assert.Equal(suite.T(), uint(2), result.Orders[1].ID)
}

func (suite *OrderPresenterTestSuite) Test_PresentOrders_WithEmptyList_ShouldReturnEmptyDTO() {
	// GIVEN an empty order list
	orders := []*entities.OrderEntity{}

	// WHEN orders are presented
	result := suite.presenter.PresentOrders(orders)

	// THEN an empty orders DTO should be returned
	assert.NotNil(suite.T(), result)
	assert.Empty(suite.T(), result.Orders)
}

// Feature: Order Presenter - Present Products
// Scenario: Enrich products with data from product service

func (suite *OrderPresenterTestSuite) Test_PresentProducts_WithEnrichedData_ShouldIncludeAllFields() {
	// GIVEN order products
	orderProducts := []*entities.OrderProductEntity{
		{ProductId: 10, Price: 25.00, Quantity: 2},
		{ProductId: 20, Price: 50.00, Quantity: 1},
	}

	enrichedProducts := []*clients.ProductDTO{
		{ID: 10, Name: "Product 1", Description: "Desc 1", Category: 1, ImageLink: "img1.jpg", Price: 25.00},
		{ID: 20, Name: "Product 2", Description: "Desc 2", Category: 2, ImageLink: "img2.jpg", Price: 50.00},
	}

	suite.mockProductClient.EXPECT().
		GetProducts(mock.Anything, []uint{10, 20}).
		Return(enrichedProducts, nil).
		Once()

	// WHEN products are presented
	result := suite.presenter.PresentProducts(orderProducts)

	// THEN all products should be enriched
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "Product 1", result[0].Name)
	assert.Equal(suite.T(), "Desc 1", result[0].Description)
	assert.Equal(suite.T(), "Product 2", result[1].Name)
	suite.mockProductClient.AssertExpectations(suite.T())
}

func (suite *OrderPresenterTestSuite) Test_PresentProducts_WithMissingProduct_ShouldReturnPartialEnrichment() {
	// GIVEN order products
	orderProducts := []*entities.OrderProductEntity{
		{ProductId: 10, Price: 25.00, Quantity: 2},
		{ProductId: 99, Price: 50.00, Quantity: 1},
	}

	// AND product service only returns one product
	enrichedProducts := []*clients.ProductDTO{
		{ID: 10, Name: "Product 1", Price: 25.00},
	}

	suite.mockProductClient.EXPECT().
		GetProducts(mock.Anything, []uint{10, 99}).
		Return(enrichedProducts, nil).
		Once()

	// WHEN products are presented
	result := suite.presenter.PresentProducts(orderProducts)

	// THEN first product should be enriched
	assert.Equal(suite.T(), "Product 1", result[0].Name)
	// AND second product should have basic data only
	assert.Equal(suite.T(), uint(99), result[1].ProductId)
	assert.Empty(suite.T(), result[1].Name)
	suite.mockProductClient.AssertExpectations(suite.T())
}

// Feature: Order Presenter - Present Status
// Scenario: Transform status entity to DTO with description

func (suite *OrderPresenterTestSuite) Test_PresentStatus_WithRecebido_ShouldIncludeDescription() {
	// GIVEN a status with value 1 (Recebido)
	now := time.Now()
	status := &entities.OrderStatusEntity{
		ID:            1,
		OrderId:       123,
		CurrentStatus: 1,
		CreatedAt:     now,
	}

	// WHEN the status is presented
	result := suite.presenter.PresentStatus(status)

	// THEN the DTO should include the description
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), uint(1), result.CurrentStatus)
	assert.Equal(suite.T(), "Recebido", result.CurrentStatusDescription)
	assert.Equal(suite.T(), uint(123), result.OrderId)
}

func (suite *OrderPresenterTestSuite) Test_PresentStatus_WithEmPreparacao_ShouldReturnCorrectDescription() {
	// GIVEN a status with value 2 (Em preparação)
	status := &entities.OrderStatusEntity{
		ID:            2,
		CurrentStatus: 2,
		OrderId:       456,
		CreatedAt:     time.Now(),
	}

	// WHEN the status is presented
	result := suite.presenter.PresentStatus(status)

	// THEN the description should be "Em preparação"
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Em preparação", result.CurrentStatusDescription)
}

func (suite *OrderPresenterTestSuite) Test_PresentStatus_WithPronto_ShouldReturnCorrectDescription() {
	// GIVEN a status with value 3 (Pronto)
	status := &entities.OrderStatusEntity{
		CurrentStatus: 3,
		OrderId:       789,
		CreatedAt:     time.Now(),
	}

	// WHEN the status is presented
	result := suite.presenter.PresentStatus(status)

	// THEN the description should be "Pronto"
	assert.Equal(suite.T(), "Pronto", result.CurrentStatusDescription)
}

func (suite *OrderPresenterTestSuite) Test_PresentStatus_WithFinalizado_ShouldReturnCorrectDescription() {
	// GIVEN a status with value 4 (Finalizado)
	status := &entities.OrderStatusEntity{
		CurrentStatus: 4,
		OrderId:       1000,
		CreatedAt:     time.Now(),
	}

	// WHEN the status is presented
	result := suite.presenter.PresentStatus(status)

	// THEN the description should be "Finalizado"
	assert.Equal(suite.T(), "Finalizado", result.CurrentStatusDescription)
}

func (suite *OrderPresenterTestSuite) Test_PresentStatus_WithInvalidStatus_ShouldReturnNil() {
	// GIVEN a status with invalid value
	status := &entities.OrderStatusEntity{
		CurrentStatus: 99,
		OrderId:       500,
		CreatedAt:     time.Now(),
	}

	// WHEN the status is presented
	result := suite.presenter.PresentStatus(status)

	// THEN nil should be returned
	assert.Nil(suite.T(), result)
}

func (suite *OrderPresenterTestSuite) Test_PresentMultipleStatus_WithMultipleStatuses_ShouldPresentAll() {
	// GIVEN multiple statuses
	now := time.Now()
	statuses := []*entities.OrderStatusEntity{
		{ID: 1, CurrentStatus: 1, OrderId: 100, CreatedAt: now},
		{ID: 2, CurrentStatus: 2, OrderId: 100, CreatedAt: now},
		{ID: 3, CurrentStatus: 3, OrderId: 100, CreatedAt: now},
	}

	// WHEN statuses are presented
	result := suite.presenter.PresentMultipleStatus(statuses)

	// THEN all statuses should be in the result
	assert.Len(suite.T(), result, 3)
	assert.Equal(suite.T(), "Recebido", result[0].CurrentStatusDescription)
	assert.Equal(suite.T(), "Em preparação", result[1].CurrentStatusDescription)
	assert.Equal(suite.T(), "Pronto", result[2].CurrentStatusDescription)
}
