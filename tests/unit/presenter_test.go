package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/viniciuscluna/tc-fiap-50/internal/infrastructure/clients"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/presenter"
	"github.com/viniciuscluna/tc-fiap-50/tests/mocks"
)

func TestOrderPresenter_Present(t *testing.T) {
	t.Run("deve apresentar pedido com dados do cliente e produtos", func(t *testing.T) {
		// Arrange
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		customerClient.EXPECT().
			GetCustomer(mock.Anything, uint(1)).
			Return(&clients.CustomerDTO{
				ID:    1,
				Name:  "João Silva",
				CPF:   12345678900,
				Email: "joao@example.com",
			}, nil)

		productClient.EXPECT().
			GetProducts(mock.Anything, []uint{101}).
			Return([]*clients.ProductDTO{
				{
					ID:          101,
					Name:        "Pizza",
					Description: "Pizza de calabresa",
					Price:       25.50,
					Category:    1,
					ImageLink:   "http://example.com/pizza.jpg",
				},
			}, nil)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		order := &entities.OrderEntity{
			ID:          1,
			CreatedAt:   time.Now(),
			TotalAmount: 51.00,
			CustomerId:  1,
			Products: []*entities.OrderProductEntity{
				{
					ID:        1,
					OrderId:   1,
					ProductId: 101,
					Price:     25.50,
					Quantity:  2,
				},
			},
			Status: []*entities.OrderStatusEntity{
				{
					ID:            1,
					CreatedAt:     time.Now(),
					CurrentStatus: 1,
					OrderId:       1,
				},
			},
		}

		// Act
		result := p.Present(order)

		// Assert
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, float32(51.00), result.TotalAmount)
		assert.NotNil(t, result.Customer)
		assert.Equal(t, "João Silva", result.Customer.Name)
		assert.Len(t, result.Products, 1)
		assert.Equal(t, "Pizza", result.Products[0].Name)
	})
}

func TestOrderPresenter_PresentProducts(t *testing.T) {
	t.Run("deve enriquecer produtos com dados do serviço", func(t *testing.T) {
		// Arrange
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		productClient.EXPECT().
			GetProducts(mock.Anything, []uint{101, 102}).
			Return([]*clients.ProductDTO{
				{ID: 101, Name: "Pizza", Description: "Deliciosa", Price: 25.50, Category: 1, ImageLink: "img1.jpg"},
				{ID: 102, Name: "Suco", Description: "Natural", Price: 8.00, Category: 2, ImageLink: "img2.jpg"},
			}, nil)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		orderProducts := []*entities.OrderProductEntity{
			{ProductId: 101, Price: 25.50, Quantity: 2},
			{ProductId: 102, Price: 8.00, Quantity: 1},
		}

		// Act
		result := p.PresentProducts(orderProducts)

		// Assert
		assert.Len(t, result, 2)
		assert.Equal(t, "Pizza", result[0].Name)
		assert.Equal(t, "Suco", result[1].Name)
		assert.Equal(t, "Deliciosa", result[0].Description)
	})

	t.Run("deve retornar produtos sem enriquecimento quando serviço falha", func(t *testing.T) {
		// Arrange
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		productClient.EXPECT().
			GetProducts(mock.Anything, mock.Anything).
			Return(nil, assert.AnError)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		orderProducts := []*entities.OrderProductEntity{
			{ProductId: 101, Price: 25.50, Quantity: 2},
		}

		// Act
		result := p.PresentProducts(orderProducts)

		// Assert
		assert.Len(t, result, 1)
		assert.Equal(t, "", result[0].Name) // Sem enriquecimento
		assert.Equal(t, uint(101), result[0].ProductId)
	})
}

func TestOrderPresenter_PresentStatus(t *testing.T) {
	t.Run("deve apresentar status com descrição", func(t *testing.T) {
		// Arrange
		customerClient := mocks.NewMockCustomerClient(t)
		productClient := mocks.NewMockProductClient(t)

		p := presenter.NewOrderPresenterImpl(customerClient, productClient)

		status := &entities.OrderStatusEntity{
			ID:            1,
			CreatedAt:     time.Now(),
			CurrentStatus: 2,
			OrderId:       1,
		}

		// Act
		result := p.PresentStatus(status)

		// Assert
		assert.NotNil(t, result)
		assert.Equal(t, uint(2), result.CurrentStatus)
		assert.Equal(t, "Em preparação", result.CurrentStatusDescription)
	})
}

func TestGetStatusDescription(t *testing.T) {
	tests := []struct {
		status      uint
		expected    string
		shouldError bool
	}{
		{1, "Recebido", false},
		{2, "Em preparação", false},
		{3, "Pronto", false},
		{4, "Finalizado", false},
		{5, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result, err := presenter.GetStatusDescription(tt.status)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
