package secondary_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	secondary "github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/persistence"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OrderProductRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository *secondary.OrderProductRepositoryImpl
}

func (suite *OrderProductRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	err = db.AutoMigrate(&entities.OrderEntity{}, &entities.OrderProductEntity{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repository = secondary.NewOrderProductRepositoryImpl(db)
}

func (suite *OrderProductRepositoryTestSuite) TearDownTest() {
	sqlDB, err := suite.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func TestOrderProductRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderProductRepositoryTestSuite))
}

// Feature: Order Product Repository - Add Order Product
// Scenario: Add a product to an order successfully

func (suite *OrderProductRepositoryTestSuite) Test_AddOrderProduct_WithValidData_ShouldCreateSuccessfully() {
	// GIVEN an existing order
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(order)

	// AND a valid order product entity
	orderProduct := &entities.OrderProductEntity{
		OrderId:   order.ID,
		ProductId: 10,
		Price:     25.50,
		Quantity:  2,
	}

	// WHEN the order product is added to the repository
	err := suite.repository.AddOrderProduct(orderProduct)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND the order product should be persisted in the database
	var savedProduct entities.OrderProductEntity
	suite.db.Where("order_id = ? AND product_id = ?", order.ID, 10).First(&savedProduct)
	assert.NotZero(suite.T(), savedProduct.ID)
	assert.Equal(suite.T(), orderProduct.OrderId, savedProduct.OrderId)
	assert.Equal(suite.T(), orderProduct.ProductId, savedProduct.ProductId)
	assert.Equal(suite.T(), orderProduct.Price, savedProduct.Price)
	assert.Equal(suite.T(), orderProduct.Quantity, savedProduct.Quantity)
}

func (suite *OrderProductRepositoryTestSuite) Test_AddOrderProduct_WithMultipleProducts_ShouldCreateAll() {
	// GIVEN an existing order
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 200.00}
	suite.db.Create(order)

	// AND multiple order product entities
	product1 := &entities.OrderProductEntity{
		OrderId:   order.ID,
		ProductId: 10,
		Price:     50.00,
		Quantity:  2,
	}
	product2 := &entities.OrderProductEntity{
		OrderId:   order.ID,
		ProductId: 20,
		Price:     100.00,
		Quantity:  1,
	}

	// WHEN both products are added
	err1 := suite.repository.AddOrderProduct(product1)
	err2 := suite.repository.AddOrderProduct(product2)

	// THEN both operations should complete without errors
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	// AND both products should be persisted
	var count int64
	suite.db.Model(&entities.OrderProductEntity{}).Where("order_id = ?", order.ID).Count(&count)
	assert.Equal(suite.T(), int64(2), count)
}

func (suite *OrderProductRepositoryTestSuite) Test_AddOrderProduct_WithDifferentPrices_ShouldCreateWithCorrectValues() {
	// GIVEN an existing order
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(order)

	// AND order products with different prices
	product1 := &entities.OrderProductEntity{
		OrderId:   order.ID,
		ProductId: 10,
		Price:     10.50,
		Quantity:  1,
	}
	product2 := &entities.OrderProductEntity{
		OrderId:   order.ID,
		ProductId: 20,
		Price:     89.50,
		Quantity:  1,
	}

	// WHEN both products are added
	err1 := suite.repository.AddOrderProduct(product1)
	err2 := suite.repository.AddOrderProduct(product2)

	// THEN both operations should complete without errors
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	// AND prices should be preserved correctly
	var saved1 entities.OrderProductEntity
	suite.db.Where("product_id = ?", 10).First(&saved1)
	assert.Equal(suite.T(), float32(10.50), saved1.Price)
}

func (suite *OrderProductRepositoryTestSuite) Test_AddOrderProduct_WithZeroQuantity_ShouldCreate() {
	// GIVEN an existing order
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 0}
	suite.db.Create(order)

	// AND an order product with zero quantity
	orderProduct := &entities.OrderProductEntity{
		OrderId:   order.ID,
		ProductId: 10,
		Price:     25.00,
		Quantity:  0,
	}

	// WHEN the order product is added
	err := suite.repository.AddOrderProduct(orderProduct)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	// AND the product should be saved with zero quantity
	var savedProduct entities.OrderProductEntity
	suite.db.Where("order_id = ? AND product_id = ?", order.ID, 10).First(&savedProduct)
	assert.Equal(suite.T(), uint(0), savedProduct.Quantity)
}

func (suite *OrderProductRepositoryTestSuite) Test_AddOrderProduct_WithSameProductMultipleTimes_ShouldCreateSeparateRecords() {
	// GIVEN an existing order
	order := &entities.OrderEntity{CustomerId: 1, TotalAmount: 100.00}
	suite.db.Create(order)

	// AND the same product added multiple times
	product1 := &entities.OrderProductEntity{
		OrderId:   order.ID,
		ProductId: 10,
		Price:     25.00,
		Quantity:  1,
	}
	product2 := &entities.OrderProductEntity{
		OrderId:   order.ID,
		ProductId: 10,
		Price:     25.00,
		Quantity:  1,
	}

	// WHEN both are added
	err1 := suite.repository.AddOrderProduct(product1)
	err2 := suite.repository.AddOrderProduct(product2)

	// THEN both operations should complete without errors
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	// AND two separate records should exist
	var count int64
	suite.db.Model(&entities.OrderProductEntity{}).Where("order_id = ? AND product_id = ?", order.ID, 10).Count(&count)
	assert.Equal(suite.T(), int64(2), count)
}
