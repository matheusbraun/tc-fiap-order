package clients

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	mockHTTPClient "github.com/viniciuscluna/tc-fiap-50/mocks/shared/httpclient"
)

// Feature: Customer Client
// Scenario: Fetch customer data from external customer service
type CustomerClientTestSuite struct {
	suite.Suite
	mockHTTPClient *mockHTTPClient.MockHTTPClient
	client         CustomerClient
	baseURL        string
}

func (suite *CustomerClientTestSuite) SetupTest() {
	suite.mockHTTPClient = mockHTTPClient.NewMockHTTPClient(suite.T())
	suite.baseURL = "http://customer-service"
	suite.client = NewCustomerClientImpl(suite.mockHTTPClient, suite.baseURL)
}

func TestCustomerClientTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerClientTestSuite))
}

// Scenario: Get customer with valid ID should return customer data
func (suite *CustomerClientTestSuite) Test_GetCustomer_WithValidId_ShouldReturnCustomer() {
	// GIVEN a valid customer ID
	customerID := uint(123)
	expectedURL := "http://customer-service/v1/customer/123"
	ctx := context.Background()

	expectedCustomer := &CustomerDTO{
		ID:    123,
		Name:  "John Doe",
		CPF:   12345678900,
		Email: "john.doe@example.com",
	}

	// AND the HTTP client returns customer data
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &CustomerDTO{}).
		RunAndReturn(func(ctx context.Context, url string, target interface{}) error {
			// Populate the target with expected data
			if dto, ok := target.(*CustomerDTO); ok {
				*dto = *expectedCustomer
			}
			return nil
		}).
		Once()

	// WHEN GetCustomer is called
	result, err := suite.client.GetCustomer(ctx, customerID)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedCustomer.ID, result.ID)
	assert.Equal(suite.T(), expectedCustomer.Name, result.Name)
	assert.Equal(suite.T(), expectedCustomer.CPF, result.CPF)
	assert.Equal(suite.T(), expectedCustomer.Email, result.Email)
}

// Scenario: Get customer with HTTP error should return error
func (suite *CustomerClientTestSuite) Test_GetCustomer_WithHTTPError_ShouldReturnError() {
	// GIVEN a customer ID
	customerID := uint(456)
	expectedURL := "http://customer-service/v1/customer/456"
	ctx := context.Background()

	// AND the HTTP client returns an error
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &CustomerDTO{}).
		Return(errors.New("connection timeout")).
		Once()

	// WHEN GetCustomer is called
	result, err := suite.client.GetCustomer(ctx, customerID)

	// THEN the operation should return an error
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to fetch customer 456")
	assert.Contains(suite.T(), err.Error(), "connection timeout")
}

// Scenario: Get customer with service unavailable should return error
func (suite *CustomerClientTestSuite) Test_GetCustomer_WithServiceUnavailable_ShouldReturnError() {
	// GIVEN a customer ID
	customerID := uint(789)
	expectedURL := "http://customer-service/v1/customer/789"
	ctx := context.Background()

	// AND the service is unavailable
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &CustomerDTO{}).
		Return(errors.New("503 service unavailable")).
		Once()

	// WHEN GetCustomer is called
	result, err := suite.client.GetCustomer(ctx, customerID)

	// THEN the operation should return an error
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to fetch customer")
}

// Scenario: Get customer with different customer IDs should construct correct URLs
func (suite *CustomerClientTestSuite) Test_GetCustomer_WithDifferentIds_ShouldConstructCorrectURLs() {
	// GIVEN multiple customer IDs
	testCases := []struct {
		customerID  uint
		expectedURL string
	}{
		{1, "http://customer-service/v1/customer/1"},
		{999, "http://customer-service/v1/customer/999"},
		{12345, "http://customer-service/v1/customer/12345"},
	}

	for _, tc := range testCases {
		// GIVEN a customer ID
		ctx := context.Background()

		// AND the HTTP client is configured for this ID
		suite.mockHTTPClient.EXPECT().
			Get(ctx, tc.expectedURL, &CustomerDTO{}).
			Return(nil).
			Once()

		// WHEN GetCustomer is called
		_, _ = suite.client.GetCustomer(ctx, tc.customerID)

		// THEN the correct URL should be constructed (verified by mock expectation)
	}
}

// Scenario: Get customer with context cancellation should propagate error
func (suite *CustomerClientTestSuite) Test_GetCustomer_WithCancelledContext_ShouldReturnError() {
	// GIVEN a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	customerID := uint(100)
	expectedURL := "http://customer-service/v1/customer/100"

	// AND the HTTP client returns context cancelled error
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &CustomerDTO{}).
		Return(context.Canceled).
		Once()

	// WHEN GetCustomer is called with cancelled context
	result, err := suite.client.GetCustomer(ctx, customerID)

	// THEN the operation should return an error
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to fetch customer")
}

// Scenario: Get customer with zero ID should still attempt fetch
func (suite *CustomerClientTestSuite) Test_GetCustomer_WithZeroId_ShouldAttemptFetch() {
	// GIVEN a zero customer ID
	customerID := uint(0)
	expectedURL := "http://customer-service/v1/customer/0"
	ctx := context.Background()

	// AND the HTTP client is configured
	suite.mockHTTPClient.EXPECT().
		Get(ctx, expectedURL, &CustomerDTO{}).
		Return(errors.New("customer not found")).
		Once()

	// WHEN GetCustomer is called
	result, err := suite.client.GetCustomer(ctx, customerID)

	// THEN the operation should return an error
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}
