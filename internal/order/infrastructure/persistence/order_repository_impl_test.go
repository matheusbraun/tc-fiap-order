package secondary_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	secondary "github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/persistence"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository *secondary.OrderRepositoryImpl
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	err = db.AutoMigrate(&entities.OrderEntity{}, &entities.OrderProductEntity{}, &entities.OrderStatusEntity{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repository = secondary.NewOrderRepositoryImpl(db)
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	sqlDB, err := suite.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func TestOrderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

// Feature: Order Repository - Add Order
// Scenario: Create a new order successfully

func (suite *OrderRepositoryTestSuite) Test_AddOrder_WithValidData_ShouldCreateSuccessfully() {
	// GIVEN a valid order entity
	order := &entities.OrderEntity{
		CustomerId:  1,
		TotalAmount: 100.50,
		CreatedAt:   time.Now(),
	}

	// WHEN the order is added to the repository
	result, err := suite.repository.AddOrder(order)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND the result should not be nil
	assert.NotNil(suite.T(), result)
	// AND the order should have an ID assigned
	assert.NotZero(suite.T(), result.ID)
	// AND all order fields should be preserved
	assert.Equal(suite.T(), order.CustomerId, result.CustomerId)
	assert.Equal(suite.T(), order.TotalAmount, result.TotalAmount)
}

func (suite *OrderRepositoryTestSuite) Test_AddOrder_WithMinimalData_ShouldCreateWithDefaults() {
	// GIVEN an order with minimal required data
	order := &entities.OrderEntity{
		CustomerId: 2,
	}

	// WHEN the order is added
	result, err := suite.repository.AddOrder(order)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// AND default values should be applied
	assert.NotZero(suite.T(), result.ID)
	assert.NotZero(suite.T(), result.CreatedAt)
}

// Feature: Order Repository - Get Order
// Scenario: Retrieve an order by ID with preloaded relations

func (suite *OrderRepositoryTestSuite) Test_GetOrder_WithValidId_ShouldReturnOrderWithPreloadedData() {
	// GIVEN an existing order with products and status
	order := &entities.OrderEntity{
		CustomerId:  1,
		TotalAmount: 150.00,
	}
	suite.db.Create(order)

	product := &entities.OrderProductEntity{
		OrderId:   order.ID,
		ProductId: 10,
		Price:     50.00,
		Quantity:  3,
	}
	suite.db.Create(product)

	status := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 1,
		CreatedAt:     time.Now(),
	}
	suite.db.Create(status)

	// WHEN the order is retrieved by ID
	result, err := suite.repository.GetOrder(order.ID)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// AND all order data should be correct
	assert.Equal(suite.T(), order.ID, result.ID)
	assert.Equal(suite.T(), order.CustomerId, result.CustomerId)
	// AND products should be preloaded
	assert.NotNil(suite.T(), result.Products)
	assert.Len(suite.T(), result.Products, 1)
	assert.Equal(suite.T(), product.ProductId, result.Products[0].ProductId)
	// AND status should be preloaded
	assert.NotNil(suite.T(), result.Status)
	assert.Len(suite.T(), result.Status, 1)
	assert.Equal(suite.T(), status.CurrentStatus, result.Status[0].CurrentStatus)
}

func (suite *OrderRepositoryTestSuite) Test_GetOrder_WithInvalidId_ShouldReturnError() {
	// GIVEN a non-existent order ID
	nonExistentId := uint(9999)

	// WHEN attempting to retrieve the order
	result, err := suite.repository.GetOrder(nonExistentId)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	// AND the result should be nil
	assert.Nil(suite.T(), result)
	// AND the error should be a record not found error
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *OrderRepositoryTestSuite) Test_GetOrder_WithMultipleStatuses_ShouldOrderByCreatedAtDesc() {
	// GIVEN an order with multiple status updates
	order := &entities.OrderEntity{
		CustomerId:  1,
		TotalAmount: 200.00,
	}
	suite.db.Create(order)

	oldStatus := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 1,
		CreatedAt:     time.Now().Add(-1 * time.Hour),
	}
	suite.db.Create(oldStatus)

	newStatus := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 2,
		CreatedAt:     time.Now(),
	}
	suite.db.Create(newStatus)

	// WHEN the order is retrieved
	result, err := suite.repository.GetOrder(order.ID)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND both statuses should be returned
	assert.Len(suite.T(), result.Status, 2)
	// AND the first status should be the most recent (DESC order)
	assert.Equal(suite.T(), newStatus.CurrentStatus, result.Status[0].CurrentStatus)
	assert.Equal(suite.T(), oldStatus.CurrentStatus, result.Status[1].CurrentStatus)
}

// Feature: Order Repository - Get Orders
// Scenario: List all orders excluding finished ones

func (suite *OrderRepositoryTestSuite) Test_GetOrders_ShouldExcludeFinishedOrders() {
	// GIVEN multiple orders with different statuses
	activeOrder1 := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(activeOrder1)
	suite.db.Create(&entities.OrderStatusEntity{OrderId: activeOrder1.ID, CurrentStatus: 1})

	activeOrder2 := &entities.OrderEntity{CustomerId: 2, TotalAmount: 150.00}
	suite.db.Create(activeOrder2)
	suite.db.Create(&entities.OrderStatusEntity{OrderId: activeOrder2.ID, CurrentStatus: 2})

	finishedOrder := &entities.OrderEntity{CustomerId: 3, TotalAmount: 200.00}
	suite.db.Create(finishedOrder)
	suite.db.Create(&entities.OrderStatusEntity{OrderId: finishedOrder.ID, CurrentStatus: 4})

	// WHEN all orders are retrieved
	results, err := suite.repository.GetOrders()

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), results)
	// AND only active orders should be returned (excluding status 4)
	assert.Len(suite.T(), results, 2)
	// AND the finished order should not be in the results
	for _, order := range results {
		assert.NotEqual(suite.T(), finishedOrder.ID, order.ID)
	}
}

func (suite *OrderRepositoryTestSuite) Test_GetOrders_ShouldOrderByCreatedAtAsc() {
	// GIVEN multiple active orders created at different times
	oldOrder := &entities.OrderEntity{
		CustomerId:  1,
		TotalAmount: 100.00,
		CreatedAt:   time.Now().Add(-2 * time.Hour),
	}
	suite.db.Create(oldOrder)
	suite.db.Create(&entities.OrderStatusEntity{OrderId: oldOrder.ID, CurrentStatus: 1})

	newOrder := &entities.OrderEntity{
		CustomerId:  2,
		TotalAmount: 150.00,
		CreatedAt:   time.Now(),
	}
	suite.db.Create(newOrder)
	suite.db.Create(&entities.OrderStatusEntity{OrderId: newOrder.ID, CurrentStatus: 1})

	// WHEN all orders are retrieved
	results, err := suite.repository.GetOrders()

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 2)
	// AND orders should be sorted by creation time ascending (oldest first)
	assert.Equal(suite.T(), oldOrder.ID, results[0].ID)
	assert.Equal(suite.T(), newOrder.ID, results[1].ID)
}

func (suite *OrderRepositoryTestSuite) Test_GetOrders_WithNoOrders_ShouldReturnEmptyList() {
	// GIVEN an empty database
	// WHEN all orders are retrieved
	results, err := suite.repository.GetOrders()

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND an empty list should be returned
	assert.NotNil(suite.T(), results)
	assert.Empty(suite.T(), results)
}

func (suite *OrderRepositoryTestSuite) Test_GetOrders_ShouldPreloadProductsAndStatus() {
	// GIVEN an order with products and status
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(order)

	product := &entities.OrderProductEntity{
		OrderId:   order.ID,
		ProductId: 5,
		Price:     25.00,
		Quantity:  4,
	}
	suite.db.Create(product)

	status := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 2,
	}
	suite.db.Create(status)

	// WHEN all orders are retrieved
	results, err := suite.repository.GetOrders()

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 1)
	// AND products should be preloaded
	assert.NotNil(suite.T(), results[0].Products)
	assert.Len(suite.T(), results[0].Products, 1)
	assert.Equal(suite.T(), product.ProductId, results[0].Products[0].ProductId)
	// AND status should be preloaded
	assert.NotNil(suite.T(), results[0].Status)
	assert.Len(suite.T(), results[0].Status, 1)
	assert.Equal(suite.T(), status.CurrentStatus, results[0].Status[0].CurrentStatus)
}
