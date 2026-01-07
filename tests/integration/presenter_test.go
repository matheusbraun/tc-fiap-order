package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/viniciuscluna/tc-fiap-50/internal/infrastructure/clients"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/presenter"
	"github.com/viniciuscluna/tc-fiap-50/tests/mocks"
)

func TestOrderPresenter_Present(t *testing.T) {
	t.Run("deve apresentar pedido com dados de cliente", func(t *testing.T) {
		// Setup
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		order := &entities.OrderEntity{
			ID:          1,
			CustomerId:  1,
			TotalAmount: 100.0,
			CreatedAt:   time.Now(),
			Products: []*entities.OrderProductEntity{
				{ID: 1, ProductId: 101, Price: 50.0, Quantity: 2},
			},
			Status: []*entities.OrderStatusEntity{
				{ID: 1, OrderId: 1, CurrentStatus: 1, CreatedAt: time.Now()},
			},
		}

		customerData := &clients.CustomerDTO{
			ID:    1,
			Name:  "Cliente Teste",
			Email: "cliente@teste.com",
			CPF:   12345678900,
		}

		productData := []*clients.ProductDTO{
			{
				ID:          101,
				Name:        "Produto Teste",
				Description: "Descrição",
				Price:       50.0,
				Category:    1,
				ImageLink:   "http://image.com",
			},
		}

		customerClient.On("GetCustomer", mock.Anything, uint(1)).Return(customerData, nil)
		productClient.On("GetProducts", mock.Anything, []uint{101}).Return(productData, nil)

		// Execute
		result := p.Present(order)

		// Assert
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, float32(100.0), result.TotalAmount)
		require.NotNil(t, result.Customer)
		assert.Equal(t, "Cliente Teste", result.Customer.Name)
		assert.Len(t, result.Products, 1)
		assert.Equal(t, "Produto Teste", result.Products[0].Name)
		customerClient.AssertExpectations(t)
		productClient.AssertExpectations(t)
	})

	t.Run("deve apresentar pedido mesmo se cliente não for encontrado", func(t *testing.T) {
		// Setup
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		order := &entities.OrderEntity{
			ID:          1,
			CustomerId:  1,
			TotalAmount: 100.0,
			CreatedAt:   time.Now(),
			Products:    []*entities.OrderProductEntity{},
			Status:      []*entities.OrderStatusEntity{},
		}

		customerClient.On("GetCustomer", mock.Anything, uint(1)).Return(nil, assert.AnError)
		productClient.On("GetProducts", mock.Anything, mock.Anything).Return([]*clients.ProductDTO{}, nil)

		// Execute (should not panic)
		result := p.Present(order)

		// Assert - graceful degradation
		require.NotNil(t, result)
		assert.Nil(t, result.Customer)
		customerClient.AssertExpectations(t)
	})

	t.Run("deve apresentar pedido sem cliente quando CustomerId é 0", func(t *testing.T) {
		// Setup
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		order := &entities.OrderEntity{
			ID:          1,
			CustomerId:  0, // Sem cliente
			TotalAmount: 100.0,
			CreatedAt:   time.Now(),
			Products:    []*entities.OrderProductEntity{},
			Status:      []*entities.OrderStatusEntity{},
		}

		productClient.On("GetProducts", mock.Anything, mock.Anything).Return([]*clients.ProductDTO{}, nil)

		// Execute
		result := p.Present(order)

		// Assert
		require.NotNil(t, result)
		assert.Nil(t, result.Customer)
		// Should NOT call customer client
		customerClient.AssertNotCalled(t, "GetCustomer", mock.Anything, mock.Anything)
	})
}

func TestOrderPresenter_PresentOrders(t *testing.T) {
	t.Run("deve apresentar lista de pedidos", func(t *testing.T) {
		// Setup
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		orders := []*entities.OrderEntity{
			{ID: 1, CustomerId: 1, TotalAmount: 100.0, CreatedAt: time.Now(), Products: []*entities.OrderProductEntity{}, Status: []*entities.OrderStatusEntity{}},
			{ID: 2, CustomerId: 2, TotalAmount: 200.0, CreatedAt: time.Now(), Products: []*entities.OrderProductEntity{}, Status: []*entities.OrderStatusEntity{}},
		}

		customerClient.On("GetCustomer", mock.Anything, uint(1)).Return(&clients.CustomerDTO{ID: 1, Name: "Cliente 1"}, nil)
		customerClient.On("GetCustomer", mock.Anything, uint(2)).Return(&clients.CustomerDTO{ID: 2, Name: "Cliente 2"}, nil)
		productClient.On("GetProducts", mock.Anything, mock.Anything).Return([]*clients.ProductDTO{}, nil)

		// Execute
		result := p.PresentOrders(orders)

		// Assert
		require.NotNil(t, result)
		assert.Len(t, result.Orders, 2)
		assert.Equal(t, uint(1), result.Orders[0].ID)
		assert.Equal(t, uint(2), result.Orders[1].ID)
	})
}

func TestOrderPresenter_PresentProducts(t *testing.T) {
	t.Run("deve enriquecer produtos com dados do serviço", func(t *testing.T) {
		// Setup
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		orderProducts := []*entities.OrderProductEntity{
			{ID: 1, ProductId: 101, Price: 50.0, Quantity: 2},
			{ID: 2, ProductId: 102, Price: 30.0, Quantity: 1},
		}

		productData := []*clients.ProductDTO{
			{ID: 101, Name: "Produto 1", Description: "Desc 1", Price: 50.0, Category: 1, ImageLink: "http://img1.com"},
			{ID: 102, Name: "Produto 2", Description: "Desc 2", Price: 30.0, Category: 2, ImageLink: "http://img2.com"},
		}

		productClient.On("GetProducts", context.Background(), []uint{101, 102}).Return(productData, nil)

		// Execute
		result := p.PresentProducts(orderProducts)

		// Assert
		require.Len(t, result, 2)
		assert.Equal(t, "Produto 1", result[0].Name)
		assert.Equal(t, "Produto 2", result[1].Name)
		assert.Equal(t, "http://img1.com", result[0].ImageLink)
		productClient.AssertExpectations(t)
	})

	t.Run("deve retornar produtos sem enriquecimento quando serviço falha", func(t *testing.T) {
		// Setup
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		orderProducts := []*entities.OrderProductEntity{
			{ID: 1, ProductId: 101, Price: 50.0, Quantity: 2},
		}

		productClient.On("GetProducts", context.Background(), []uint{101}).Return(nil, assert.AnError)

		// Execute
		result := p.PresentProducts(orderProducts)

		// Assert
		require.Len(t, result, 1)
		assert.Equal(t, uint(101), result[0].ProductId)
		assert.Empty(t, result[0].Name) // No enrichment
		productClient.AssertExpectations(t)
	})
}

func TestOrderPresenter_PresentStatus(t *testing.T) {
	t.Run("deve apresentar status do pedido", func(t *testing.T) {
		// Setup
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		now := time.Now()
		status := &entities.OrderStatusEntity{
			ID:            1,
			OrderId:       1,
			CurrentStatus: 1,
			CreatedAt:     now,
		}

		// Execute
		result := p.PresentStatus(status)

		// Assert
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, uint(1), result.CurrentStatus)
		assert.Equal(t, "Recebido", result.CurrentStatusDescription)
	})

	t.Run("deve retornar todos os status", func(t *testing.T) {
		tests := []struct {
			statusID    uint
			description string
		}{
			{1, "Recebido"},
			{2, "Em preparação"},
			{3, "Pronto"},
			{4, "Finalizado"},
		}

		for _, tt := range tests {
			// Setup
			customerClient := mocks.NewMockCustomerClient(t)
			productClient := mocks.NewMockProductClient(t)
			p := presenter.NewOrderPresenterImpl(customerClient, productClient)

			status := &entities.OrderStatusEntity{
				ID:            1,
				OrderId:       1,
				CurrentStatus: tt.statusID,
				CreatedAt:     time.Now(),
			}

			// Execute
			result := p.PresentStatus(status)

			// Assert
			require.NotNil(t, result)
			assert.Equal(t, tt.description, result.CurrentStatusDescription)
		}
	})

	t.Run("deve retornar nil para status inválido", func(t *testing.T) {
		// Setup
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)
		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		status := &entities.OrderStatusEntity{
			ID:            1,
			OrderId:       1,
			CurrentStatus: 999, // Status inválido
			CreatedAt:     time.Now(),
		}

		// Execute
		result := p.PresentStatus(status)

		// Assert
		assert.Nil(t, result)
	})
}

func TestOrderPresenter_PresentMultipleStatus(t *testing.T) {
	t.Run("deve apresentar múltiplos status", func(t *testing.T) {
		// Setup
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)
		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		statusList := []*entities.OrderStatusEntity{
			{ID: 1, OrderId: 1, CurrentStatus: 1, CreatedAt: time.Now()},
			{ID: 2, OrderId: 1, CurrentStatus: 2, CreatedAt: time.Now()},
		}

		// Execute
		result := p.PresentMultipleStatus(statusList)

		// Assert
		require.Len(t, result, 2)
		assert.Equal(t, "Recebido", result[0].CurrentStatusDescription)
		assert.Equal(t, "Em preparação", result[1].CurrentStatusDescription)
	})
}

func TestGetStatusDescription(t *testing.T) {
	tests := []struct {
		name        string
		status      uint
		expected    string
		expectError bool
	}{
		{"Recebido", 1, "Recebido", false},
		{"Em preparação", 2, "Em preparação", false},
		{"Pronto", 3, "Pronto", false},
		{"Finalizado", 4, "Finalizado", false},
		{"Status inválido", 999, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := presenter.GetStatusDescription(tt.status)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
