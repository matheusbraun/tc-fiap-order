package clients

import (
	"context"
	"fmt"

	"github.com/viniciuscluna/tc-fiap-50/internal/shared/httpclient"
)

type customerClientImpl struct {
	httpClient httpclient.HTTPClient
	baseURL    string
}

func NewCustomerClientImpl(httpClient httpclient.HTTPClient, baseURL string) CustomerClient {
	return &customerClientImpl{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *customerClientImpl) GetCustomer(ctx context.Context, customerID uint) (*CustomerDTO, error) {
	url := fmt.Sprintf("%s/v1/customer/%d", c.baseURL, customerID)
	var customer CustomerDTO

	if err := c.httpClient.Get(ctx, url, &customer); err != nil {
		return nil, fmt.Errorf("failed to fetch customer %d: %w", customerID, err)
	}

	return &customer, nil
}
