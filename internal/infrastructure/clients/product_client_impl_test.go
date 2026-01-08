package clients

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	mockHTTPClient "github.com/viniciuscluna/tc-fiap-50/mocks/shared/httpclient"
)

// Feature: Product Client
// Scenario: Fetch product data from external product service
type ProductClientTestSuite struct {
	suite.Suite
	mockHTTPClient *mockHTTPClient.MockHTTPClient
	client         ProductClient
	baseURL        string
}

func (suite *ProductClientTestSuite) SetupTest() {
	suite.mockHTTPClient = mockHTTPClient.NewMockHTTPClient(suite.T())
	suite.baseURL = "http://product-service"
	suite.client = NewProductClientImpl(suite.mockHTTPClient, suite.baseURL)
}

func TestProductClientTestSuite(t *testing.T) {
	suite.Run(t, new(ProductClientTestSuite))
}

// Scenario: Get product with valid ID should return product data
func (suite *ProductClientTestSuite) Test_GetProduct_WithValidId_ShouldReturnProduct() {
	// GIVEN a valid product ID
	productID := uint(456)
	expectedURL := "http://product-service/v1/product/456"
	ctx := context.Background()

	expectedProduct := &ProductDTO{
		ID:          456,
		Name:        "Hamburger",
		Description: "Delicious hamburger",
		Price:       25.90,
		Category:    1,
		ImageLink:   "http://example.com/hamburger.jpg",
	}

	// AND the HTTP client returns product data
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &ProductDTO{}).
		RunAndReturn(func(ctx context.Context, url string, target interface{}) error {
			// Populate the target with expected data
			if dto, ok := target.(*ProductDTO); ok {
				*dto = *expectedProduct
			}
			return nil
		}).
		Once()

	// WHEN GetProduct is called
	result, err := suite.client.GetProduct(ctx, productID)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedProduct.ID, result.ID)
	assert.Equal(suite.T(), expectedProduct.Name, result.Name)
	assert.Equal(suite.T(), expectedProduct.Description, result.Description)
	assert.Equal(suite.T(), expectedProduct.Price, result.Price)
	assert.Equal(suite.T(), expectedProduct.Category, result.Category)
	assert.Equal(suite.T(), expectedProduct.ImageLink, result.ImageLink)
}

// Scenario: Get product with HTTP error should return error
func (suite *ProductClientTestSuite) Test_GetProduct_WithHTTPError_ShouldReturnError() {
	// GIVEN a product ID
	productID := uint(789)
	expectedURL := "http://product-service/v1/product/789"
	ctx := context.Background()

	// AND the HTTP client returns an error
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &ProductDTO{}).
		Return(errors.New("network error")).
		Once()

	// WHEN GetProduct is called
	result, err := suite.client.GetProduct(ctx, productID)

	// THEN the operation should return an error
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to fetch product 789")
	assert.Contains(suite.T(), err.Error(), "network error")
}

// Scenario: Get products with empty list should return empty slice
func (suite *ProductClientTestSuite) Test_GetProducts_WithEmptyList_ShouldReturnEmptySlice() {
	// GIVEN an empty list of product IDs
	productIDs := []uint{}
	ctx := context.Background()

	// WHEN GetProducts is called
	result, err := suite.client.GetProducts(ctx, productIDs)

	// THEN the operation should return an empty slice without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Empty(suite.T(), result)
}

// Scenario: Get products with single ID should return single product
func (suite *ProductClientTestSuite) Test_GetProducts_WithSingleId_ShouldReturnSingleProduct() {
	// GIVEN a single product ID
	productIDs := []uint{100}
	expectedURL := "http://product-service/v1/product/100"
	ctx := context.Background()

	expectedProduct := &ProductDTO{
		ID:          100,
		Name:        "Pizza",
		Description: "Margherita pizza",
		Price:       35.00,
		Category:    1,
		ImageLink:   "http://example.com/pizza.jpg",
	}

	// AND the HTTP client returns product data
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &ProductDTO{}).
		RunAndReturn(func(ctx context.Context, url string, target interface{}) error {
			if dto, ok := target.(*ProductDTO); ok {
				*dto = *expectedProduct
			}
			return nil
		}).
		Once()

	// WHEN GetProducts is called
	result, err := suite.client.GetProducts(ctx, productIDs)

	// THEN the operation should return a single product
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), expectedProduct.ID, result[0].ID)
	assert.Equal(suite.T(), expectedProduct.Name, result[0].Name)
}

// Scenario: Get products with multiple IDs should return multiple products
func (suite *ProductClientTestSuite) Test_GetProducts_WithMultipleIds_ShouldReturnMultipleProducts() {
	// GIVEN multiple product IDs
	productIDs := []uint{1, 2, 3}
	ctx := context.Background()

	product1 := &ProductDTO{ID: 1, Name: "Product 1", Price: 10.00, Category: 1}
	product2 := &ProductDTO{ID: 2, Name: "Product 2", Price: 20.00, Category: 2}
	product3 := &ProductDTO{ID: 3, Name: "Product 3", Price: 30.00, Category: 3}

	// AND the HTTP client returns data for each product
	suite.mockHTTPClient.EXPECT().
		Get(ctx, "http://product-service/v1/product/1", &ProductDTO{}).
		RunAndReturn(func(ctx context.Context, url string, target interface{}) error {
			if dto, ok := target.(*ProductDTO); ok {
				*dto = *product1
			}
			return nil
		}).
		Once()

	suite.mockHTTPClient.EXPECT().
		Get(ctx, "http://product-service/v1/product/2", &ProductDTO{}).
		RunAndReturn(func(ctx context.Context, url string, target interface{}) error {
			if dto, ok := target.(*ProductDTO); ok {
				*dto = *product2
			}
			return nil
		}).
		Once()

	suite.mockHTTPClient.EXPECT().
		Get(ctx, "http://product-service/v1/product/3", &ProductDTO{}).
		RunAndReturn(func(ctx context.Context, url string, target interface{}) error {
			if dto, ok := target.(*ProductDTO); ok {
				*dto = *product3
			}
			return nil
		}).
		Once()

	// WHEN GetProducts is called
	result, err := suite.client.GetProducts(ctx, productIDs)

	// THEN the operation should return all products
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 3)
	assert.Equal(suite.T(), uint(1), result[0].ID)
	assert.Equal(suite.T(), uint(2), result[1].ID)
	assert.Equal(suite.T(), uint(3), result[2].ID)
}

// Scenario: Get products with one failing should return error
func (suite *ProductClientTestSuite) Test_GetProducts_WithOneFailingRequest_ShouldReturnError() {
	// GIVEN multiple product IDs
	productIDs := []uint{10, 20, 30}
	ctx := context.Background()

	product1 := &ProductDTO{ID: 10, Name: "Product 10", Price: 15.00}

	// AND the first product succeeds
	suite.mockHTTPClient.EXPECT().
		Get(ctx, "http://product-service/v1/product/10", &ProductDTO{}).
		RunAndReturn(func(ctx context.Context, url string, target interface{}) error {
			if dto, ok := target.(*ProductDTO); ok {
				*dto = *product1
			}
			return nil
		}).
		Once()

	// AND the second product fails
	suite.mockHTTPClient.EXPECT().
		Get(ctx, "http://product-service/v1/product/20", &ProductDTO{}).
		Return(errors.New("product not found")).
		Once()

	// WHEN GetProducts is called
	result, err := suite.client.GetProducts(ctx, productIDs)

	// THEN the operation should return an error
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to fetch product 20")
}

// Scenario: Get products should maintain order of requested IDs
func (suite *ProductClientTestSuite) Test_GetProducts_ShouldMaintainOrderOfRequestedIds() {
	// GIVEN product IDs in specific order
	productIDs := []uint{5, 3, 7}
	ctx := context.Background()

	product5 := &ProductDTO{ID: 5, Name: "Product 5"}
	product3 := &ProductDTO{ID: 3, Name: "Product 3"}
	product7 := &ProductDTO{ID: 7, Name: "Product 7"}

	// AND the HTTP client returns products
	suite.mockHTTPClient.EXPECT().
		Get(ctx, "http://product-service/v1/product/5", &ProductDTO{}).
		RunAndReturn(func(ctx context.Context, url string, target interface{}) error {
			if dto, ok := target.(*ProductDTO); ok {
				*dto = *product5
			}
			return nil
		}).
		Once()

	suite.mockHTTPClient.EXPECT().
		Get(ctx, "http://product-service/v1/product/3", &ProductDTO{}).
		RunAndReturn(func(ctx context.Context, url string, target interface{}) error {
			if dto, ok := target.(*ProductDTO); ok {
				*dto = *product3
			}
			return nil
		}).
		Once()

	suite.mockHTTPClient.EXPECT().
		Get(ctx, "http://product-service/v1/product/7", &ProductDTO{}).
		RunAndReturn(func(ctx context.Context, url string, target interface{}) error {
			if dto, ok := target.(*ProductDTO); ok {
				*dto = *product7
			}
			return nil
		}).
		Once()

	// WHEN GetProducts is called
	result, err := suite.client.GetProducts(ctx, productIDs)

	// THEN the products should be in the same order as requested
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3)
	assert.Equal(suite.T(), uint(5), result[0].ID)
	assert.Equal(suite.T(), uint(3), result[1].ID)
	assert.Equal(suite.T(), uint(7), result[2].ID)
}

// Scenario: Get product with service unavailable should return error
func (suite *ProductClientTestSuite) Test_GetProduct_WithServiceUnavailable_ShouldReturnError() {
	// GIVEN a product ID
	productID := uint(999)
	expectedURL := "http://product-service/v1/product/999"
	ctx := context.Background()

	// AND the service is unavailable
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &ProductDTO{}).
		Return(errors.New("503 service unavailable")).
		Once()

	// WHEN GetProduct is called
	result, err := suite.client.GetProduct(ctx, productID)

	// THEN the operation should return an error
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to fetch product")
}

// Scenario: Get product with context cancellation should propagate error
func (suite *ProductClientTestSuite) Test_GetProduct_WithCancelledContext_ShouldReturnError() {
	// GIVEN a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	productID := uint(200)
	expectedURL := "http://product-service/v1/product/200"

	// AND the HTTP client returns context cancelled error
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &ProductDTO{}).
		Return(context.Canceled).
		Once()

	// WHEN GetProduct is called with cancelled context
	result, err := suite.client.GetProduct(ctx, productID)

	// THEN the operation should return an error
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to fetch product")
}

// Scenario: Get products with context cancellation should propagate error
func (suite *ProductClientTestSuite) Test_GetProducts_WithCancelledContext_ShouldReturnError() {
	// GIVEN a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	productIDs := []uint{1, 2}
	expectedURL := "http://product-service/v1/product/1"

	// AND the HTTP client returns context cancelled error on first request
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &ProductDTO{}).
		Return(context.Canceled).
		Once()

	// WHEN GetProducts is called with cancelled context
	result, err := suite.client.GetProducts(ctx, productIDs)

	// THEN the operation should return an error
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to fetch product")
}
