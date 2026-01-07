package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viniciuscluna/tc-fiap-50/internal/infrastructure/clients"
	"github.com/viniciuscluna/tc-fiap-50/internal/shared/httpclient"
)

func TestCustomerClient_GetCustomer(t *testing.T) {
	t.Run("deve buscar cliente com sucesso", func(t *testing.T) {
		// Arrange
		expected := &clients.CustomerDTO{
			ID:    1,
			Name:  "João Silva",
			CPF:   12345678900,
			Email: "joao@example.com",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/customer/1", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expected)
		}))
		defer server.Close()

		httpClient := httpclient.NewHTTPClient(30*time.Second, 3, 100*time.Millisecond)
		client := clients.NewCustomerClientImpl(httpClient, server.URL)

		// Act
		result, err := client.GetCustomer(context.Background(), 1)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Name, result.Name)
		assert.Equal(t, expected.Email, result.Email)
	})

	t.Run("deve retornar erro quando cliente não existe", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("customer not found"))
		}))
		defer server.Close()

		httpClient := httpclient.NewHTTPClient(30*time.Second, 3, 100*time.Millisecond)
		client := clients.NewCustomerClientImpl(httpClient, server.URL)

		// Act
		_, err := client.GetCustomer(context.Background(), 999)

		// Assert
		assert.Error(t, err)
	})
}

func TestProductClient_GetProduct(t *testing.T) {
	t.Run("deve buscar produto com sucesso", func(t *testing.T) {
		// Arrange
		expected := &clients.ProductDTO{
			ID:          101,
			Name:        "Pizza",
			Description: "Deliciosa pizza de calabresa",
			Price:       25.50,
			Category:    1,
			ImageLink:   "http://example.com/pizza.jpg",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/product/101", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expected)
		}))
		defer server.Close()

		httpClient := httpclient.NewHTTPClient(30*time.Second, 3, 100*time.Millisecond)
		client := clients.NewProductClientImpl(httpClient, server.URL)

		// Act
		result, err := client.GetProduct(context.Background(), 101)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Name, result.Name)
		assert.Equal(t, expected.Price, result.Price)
	})
}

func TestProductClient_GetProducts(t *testing.T) {
	t.Run("deve buscar múltiplos produtos com sucesso", func(t *testing.T) {
		// Arrange
		products := []*clients.ProductDTO{
			{ID: 101, Name: "Pizza", Price: 25.50, Category: 1},
			{ID: 102, Name: "Suco", Price: 8.00, Category: 2},
		}

		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(products[callCount])
			callCount++
		}))
		defer server.Close()

		httpClient := httpclient.NewHTTPClient(30*time.Second, 3, 100*time.Millisecond)
		client := clients.NewProductClientImpl(httpClient, server.URL)

		// Act
		result, err := client.GetProducts(context.Background(), []uint{101, 102})

		// Assert
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Pizza", result[0].Name)
		assert.Equal(t, "Suco", result[1].Name)
	})

	t.Run("deve retornar lista vazia quando nenhum ID fornecido", func(t *testing.T) {
		// Arrange
		httpClient := httpclient.NewHTTPClient(30*time.Second, 3, 100*time.Millisecond)
		client := clients.NewProductClientImpl(httpClient, "http://localhost")

		// Act
		result, err := client.GetProducts(context.Background(), []uint{})

		// Assert
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestProductClient_ValidateProduct(t *testing.T) {
	t.Run("deve retornar true quando produto existe", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			product := &clients.ProductDTO{ID: 101, Name: "Pizza", Price: 25.50, Category: 1}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(product)
		}))
		defer server.Close()

		httpClient := httpclient.NewHTTPClient(30*time.Second, 3, 100*time.Millisecond)
		client := clients.NewProductClientImpl(httpClient, server.URL)

		// Act
		valid, err := client.ValidateProduct(context.Background(), 101)

		// Assert
		assert.NoError(t, err)
		assert.True(t, valid)
	})
}
