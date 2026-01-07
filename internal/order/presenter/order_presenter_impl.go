package presenter

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/viniciuscluna/tc-fiap-50/internal/infrastructure/clients"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"
)

var (
	_ OrderPresenter = (*OrderPresenterImpl)(nil)
)

type OrderPresenterImpl struct {
	customerClient clients.CustomerClient
	productClient  clients.ProductClient
}

func NewOrderPresenterImpl(customerClient clients.CustomerClient, productClient clients.ProductClient) *OrderPresenterImpl {
	return &OrderPresenterImpl{
		customerClient: customerClient,
		productClient:  productClient,
	}
}

func (p *OrderPresenterImpl) Present(order *entities.OrderEntity) *dto.GetOrderResponseDto {
	ctx := context.Background()

	var customer *dto.CustomerDto
	if order.CustomerId != 0 {
		// Fetch customer data from customer service
		customerData, err := p.customerClient.GetCustomer(ctx, order.CustomerId)
		if err != nil {
			// Log error but don't fail - graceful degradation
			log.Printf("failed to fetch customer %d: %v", order.CustomerId, err)
		} else {
			customer = &dto.CustomerDto{
				ID:    customerData.ID,
				Name:  customerData.Name,
				Email: customerData.Email,
				CPF:   customerData.CPF,
			}
		}
	}

	response := &dto.GetOrderResponseDto{
		ID:          order.ID,
		CreatedAt:   order.CreatedAt,
		TotalAmount: order.TotalAmount,
		CustomerId:  order.CustomerId,
		Customer:    customer,
		Products:    p.PresentProducts(order.Products),
		Status:      p.PresentMultipleStatus(order.Status),
	}

	return response
}

func (p *OrderPresenterImpl) PresentOrders(orders []*entities.OrderEntity) *dto.GetOrdersResponseDto {
	orderDto := make([]*dto.GetOrderResponseDto, len(orders))

	for i, order := range orders {
		orderDto[i] = p.Present(order)
	}

	return &dto.GetOrdersResponseDto{
		Orders: orderDto,
	}
}

func (p *OrderPresenterImpl) PresentProducts(orderProducts []*entities.OrderProductEntity) []*dto.OrderProductDto {
	ctx := context.Background()
	orderProductDtoArr := make([]*dto.OrderProductDto, len(orderProducts))

	// Collect all product IDs
	productIDs := make([]uint, len(orderProducts))
	for i, orderProduct := range orderProducts {
		productIDs[i] = orderProduct.ProductId
	}

	// Fetch all products in batch from product service
	products, err := p.productClient.GetProducts(ctx, productIDs)
	if err != nil {
		log.Printf("failed to fetch products: %v", err)
		// Return products without enriched data
		for i, orderProduct := range orderProducts {
			orderProductDtoArr[i] = &dto.OrderProductDto{
				ProductId: orderProduct.ProductId,
				Price:     orderProduct.Price,
				Quantity:  orderProduct.Quantity,
			}
		}
		return orderProductDtoArr
	}

	// Create a map for quick lookup
	productMap := make(map[uint]*clients.ProductDTO)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// Enrich order products with product data
	for i, orderProduct := range orderProducts {
		product, exists := productMap[orderProduct.ProductId]
		if exists {
			orderProductDtoArr[i] = &dto.OrderProductDto{
				ProductId:   orderProduct.ProductId,
				Price:       orderProduct.Price,
				Quantity:    orderProduct.Quantity,
				Name:        product.Name,
				ImageLink:   product.ImageLink,
				Description: product.Description,
				Category:    product.Category,
			}
		} else {
			// Product not found, return without enriched data
			orderProductDtoArr[i] = &dto.OrderProductDto{
				ProductId: orderProduct.ProductId,
				Price:     orderProduct.Price,
				Quantity:  orderProduct.Quantity,
			}
		}
	}

	return orderProductDtoArr
}

func (p *OrderPresenterImpl) PresentStatus(orderStatus *entities.OrderStatusEntity) *dto.GetOrderStatusResponseDto {
	statusDescription, err := GetStatusDescription(orderStatus.CurrentStatus)
	if err != nil {
		return nil
	}

	return &dto.GetOrderStatusResponseDto{
		ID:                       orderStatus.ID,
		CreatedAt:                orderStatus.CreatedAt.Format(time.RFC3339),
		CurrentStatus:            orderStatus.CurrentStatus,
		CurrentStatusDescription: statusDescription,
		OrderId:                  orderStatus.OrderId,
	}
}

func (p *OrderPresenterImpl) PresentMultipleStatus(orderStatus []*entities.OrderStatusEntity) []*dto.GetOrderStatusResponseDto {
	orderStatusDtoArr := make([]*dto.GetOrderStatusResponseDto, len(orderStatus))

	for i, currentOrderStatus := range orderStatus {
		orderStatusDtoArr[i] = p.PresentStatus(currentOrderStatus)
	}

	return orderStatusDtoArr
}

// Obtain Description from CurrentStatus (id)
// 1 - Recebido
// 2 - Em preparação
// 3 - Pronto
// 4 - Finalizado
func GetStatusDescription(status uint) (string, error) {
	switch status {
	case 1:
		return "Recebido", nil
	case 2:
		return "Em preparação", nil
	case 3:
		return "Pronto", nil
	case 4:
		return "Finalizado", nil
	default:
		return "", errors.New("status not found")
	}
}
