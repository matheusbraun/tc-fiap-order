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

type OrderStatusRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository *secondary.OrderStatusRepositoryImpl
}

func (suite *OrderStatusRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	err = db.AutoMigrate(&entities.OrderEntity{}, &entities.OrderStatusEntity{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repository = secondary.NewOrderStatusRepositoryImpl(db)
}

func (suite *OrderStatusRepositoryTestSuite) TearDownTest() {
	sqlDB, err := suite.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func TestOrderStatusRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderStatusRepositoryTestSuite))
}

// Feature: Order Status Repository - Add Order Status
// Scenario: Add a status to an order successfully

func (suite *OrderStatusRepositoryTestSuite) Test_AddOrderStatus_WithValidData_ShouldCreateSuccessfully() {
	// GIVEN an existing order
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(order)

	// AND a valid order status entity
	orderStatus := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 1,
		CreatedAt:     time.Now(),
	}

	// WHEN the order status is added to the repository
	err := suite.repository.AddOrderStatus(orderStatus)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND the order status should be persisted in the database
	var savedStatus entities.OrderStatusEntity
	suite.db.Where("order_id = ?", order.ID).First(&savedStatus)
	assert.NotZero(suite.T(), savedStatus.ID)
	assert.Equal(suite.T(), orderStatus.OrderId, savedStatus.OrderId)
	assert.Equal(suite.T(), orderStatus.CurrentStatus, savedStatus.CurrentStatus)
}

func (suite *OrderStatusRepositoryTestSuite) Test_AddOrderStatus_WithStatusRecebido_ShouldCreateWithStatus1() {
	// GIVEN an existing order
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 50.00}
	suite.db.Create(order)

	// AND an order status for "Recebido" (status 1)
	orderStatus := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 1,
	}

	// WHEN the status is added
	err := suite.repository.AddOrderStatus(orderStatus)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND the status should be set to 1 (Recebido)
	var savedStatus entities.OrderStatusEntity
	suite.db.Where("order_id = ?", order.ID).First(&savedStatus)
	assert.Equal(suite.T(), uint(1), savedStatus.CurrentStatus)
}

func (suite *OrderStatusRepositoryTestSuite) Test_AddOrderStatus_WithMultipleStatuses_ShouldCreateHistory() {
	// GIVEN an existing order
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(order)

	// AND multiple status updates
	status1 := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 1,
		CreatedAt:     time.Now().Add(-2 * time.Hour),
	}
	status2 := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 2,
		CreatedAt:     time.Now().Add(-1 * time.Hour),
	}
	status3 := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 3,
		CreatedAt:     time.Now(),
	}

	// WHEN all statuses are added
	err1 := suite.repository.AddOrderStatus(status1)
	err2 := suite.repository.AddOrderStatus(status2)
	err3 := suite.repository.AddOrderStatus(status3)

	// THEN all operations should complete without errors
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NoError(suite.T(), err3)
	// AND all three status records should exist
	var count int64
	suite.db.Model(&entities.OrderStatusEntity{}).Where("order_id = ?", order.ID).Count(&count)
	assert.Equal(suite.T(), int64(3), count)
}

func (suite *OrderStatusRepositoryTestSuite) Test_AddOrderStatus_WithDifferentStatuses_ShouldCreateCorrectly() {
	// GIVEN an existing order
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(order)

	// AND status updates with different values
	status1 := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 1,
	}
	status2 := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 3,
	}

	// WHEN both statuses are added
	err1 := suite.repository.AddOrderStatus(status1)
	err2 := suite.repository.AddOrderStatus(status2)

	// THEN both operations should complete without errors
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	// AND both statuses should be stored correctly
	var count int64
	suite.db.Model(&entities.OrderStatusEntity{}).Where("order_id = ?", order.ID).Count(&count)
	assert.Equal(suite.T(), int64(2), count)
}

// Feature: Order Status Repository - Get Order Status
// Scenario: Retrieve the latest status for an order

func (suite *OrderStatusRepositoryTestSuite) Test_GetOrderStatus_WithValidOrderId_ShouldReturnLatestStatus() {
	// GIVEN an existing order with multiple statuses
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(order)

	oldStatus := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 1,
		CreatedAt:     time.Now().Add(-1 * time.Hour),
	}
	suite.db.Create(oldStatus)

	latestStatus := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 2,
		CreatedAt:     time.Now(),
	}
	suite.db.Create(latestStatus)

	// WHEN the order status is retrieved
	result, err := suite.repository.GetOrderStatus(order.ID)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// AND the latest status should be returned
	assert.Equal(suite.T(), latestStatus.ID, result.ID)
	assert.Equal(suite.T(), latestStatus.CurrentStatus, result.CurrentStatus)
	assert.Equal(suite.T(), latestStatus.OrderId, result.OrderId)
}

func (suite *OrderStatusRepositoryTestSuite) Test_GetOrderStatus_WithSingleStatus_ShouldReturnThatStatus() {
	// GIVEN an order with a single status
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(order)

	status := &entities.OrderStatusEntity{
		OrderId:       order.ID,
		CurrentStatus: 1,
	}
	suite.db.Create(status)

	// WHEN the order status is retrieved
	result, err := suite.repository.GetOrderStatus(order.ID)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// AND the single status should be returned
	assert.Equal(suite.T(), status.ID, result.ID)
	assert.Equal(suite.T(), status.CurrentStatus, result.CurrentStatus)
}

func (suite *OrderStatusRepositoryTestSuite) Test_GetOrderStatus_WithInvalidOrderId_ShouldReturnError() {
	// GIVEN a non-existent order ID
	nonExistentId := uint(9999)

	// WHEN attempting to retrieve the order status
	result, err := suite.repository.GetOrderStatus(nonExistentId)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	// AND the result should be nil
	assert.Nil(suite.T(), result)
	// AND the error should be a record not found error
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *OrderStatusRepositoryTestSuite) Test_GetOrderStatus_WithNoStatus_ShouldReturnError() {
	// GIVEN an order without any status
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(order)

	// WHEN attempting to retrieve the order status
	result, err := suite.repository.GetOrderStatus(order.ID)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *OrderStatusRepositoryTestSuite) Test_GetOrderStatus_WithAllStatusTypes_ShouldReturnLatest() {
	// GIVEN an order that went through all status stages
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 150.00}
	suite.db.Create(order)

	statuses := []uint{1, 2, 3, 4}
	for i, statusValue := range statuses {
		status := &entities.OrderStatusEntity{
			OrderId:       order.ID,
			CurrentStatus: statusValue,
			CreatedAt:     time.Now().Add(time.Duration(i) * time.Hour),
		}
		suite.db.Create(status)
	}

	// WHEN the order status is retrieved
	result, err := suite.repository.GetOrderStatus(order.ID)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND the latest status (Finalizado - 4) should be returned
	assert.Equal(suite.T(), uint(4), result.CurrentStatus)
}
