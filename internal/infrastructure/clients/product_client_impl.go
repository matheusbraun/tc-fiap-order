package clients

import (
	"context"
	"fmt"
	"strings"

	"github.com/viniciuscluna/tc-fiap-50/internal/shared/httpclient"
)

type productClientImpl struct {
	httpClient httpclient.HTTPClient
	baseURL    string
}

func NewProductClientImpl(httpClient httpclient.HTTPClient, baseURL string) ProductClient {
	return &productClientImpl{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *productClientImpl) GetProduct(ctx context.Context, productID uint) (*ProductDTO, error) {
	url := fmt.Sprintf("%s/v1/product/%d", c.baseURL, productID)
	var product ProductDTO

	if err := c.httpClient.Get(ctx, url, &product); err != nil {
		return nil, fmt.Errorf("failed to fetch product %d: %w", productID, err)
	}

	return &product, nil
}

func (c *productClientImpl) GetProducts(ctx context.Context, productIDs []uint) ([]*ProductDTO, error) {
	if len(productIDs) == 0 {
		return []*ProductDTO{}, nil
	}

	// Fetch products individually (can be optimized later with batch endpoint)
	products := make([]*ProductDTO, 0, len(productIDs))
	for _, productID := range productIDs {
		product, err := c.GetProduct(ctx, productID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product %d: %w", productID, err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (c *productClientImpl) ValidateProduct(ctx context.Context, productID uint) (bool, error) {
	product, err := c.GetProduct(ctx, productID)
	if err != nil {
		// Check if it's a 404 error
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}

	return product != nil && product.ID == productID, nil
}
