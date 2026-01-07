package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/mock"
	"github.com/viniciuscluna/tc-fiap-50/internal/infrastructure/clients"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"
	addorder "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/addOrder"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/commands"
	getorder "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrder"
	getorders "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrders"
	updateorderstatus "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/updateOrderStatus"
	"github.com/viniciuscluna/tc-fiap-50/tests/mocks"
)

type testContext struct {
	// Mocks
	orderRepo        *mocks.MockOrderRepository
	orderProductRepo *mocks.MockOrderProductRepository
	orderStatusRepo  *mocks.MockOrderStatusRepository
	customerClient   *mocks.MockCustomerClient
	productClient    *mocks.MockProductClient

	// Use cases
	addOrderUseCase     addorder.AddOrderUseCase
	getOrderUseCase     getorder.GetOrderUseCase
	getOrdersUseCase    getorders.GetOrdersUseCase
	updateStatusUseCase updateorderstatus.UpdateOrderStatusUseCase

	// Test data
	products      map[uint]*clients.ProductDTO
	customers     map[uint]*clients.CustomerDTO
	orders        map[uint]*entities.OrderEntity
	orderProducts []*dto.AddOrderProductDto
	totalAmount   float32

	// Results
	resultOrderID string
	resultOrder   *entities.OrderEntity
	resultOrders  []*entities.OrderEntity
	resultError   error
}

func (tc *testContext) reset() {
	tc.orderRepo = mocks.NewMockOrderRepository(nil)
	tc.orderProductRepo = mocks.NewMockOrderProductRepository(nil)
	tc.orderStatusRepo = mocks.NewMockOrderStatusRepository(nil)
	tc.customerClient = mocks.NewMockCustomerClient(nil)
	tc.productClient = mocks.NewMockProductClient(nil)

	tc.addOrderUseCase = addorder.NewAddOrderUseCaseImpl(
		tc.orderRepo,
		tc.orderProductRepo,
		tc.orderStatusRepo,
		tc.customerClient,
		tc.productClient,
	)

	tc.getOrderUseCase = getorder.NewGetOrderUseCaseImpl(tc.orderRepo)
	tc.getOrdersUseCase = getorders.NewGetOrdersUseCaseImpl(tc.orderRepo)
	tc.updateStatusUseCase = updateorderstatus.NewUpdateOrderStatusUseCaseImpl(tc.orderStatusRepo)

	tc.products = make(map[uint]*clients.ProductDTO)
	tc.customers = make(map[uint]*clients.CustomerDTO)
	tc.orders = make(map[uint]*entities.OrderEntity)
	tc.orderProducts = []*dto.AddOrderProductDto{}
	tc.totalAmount = 0
	tc.resultOrderID = ""
	tc.resultOrder = nil
	tc.resultOrders = nil
	tc.resultError = nil
}

// Step definitions
func (tc *testContext) oServicoDeClientesEstaDisponivel() error {
	// Mock já está pronto, nada a fazer
	return nil
}

func (tc *testContext) oServicoDeProdutosEstaDisponivel() error {
	// Mock já está pronto, nada a fazer
	return nil
}

func (tc *testContext) oClienteComIDExiste(customerID int) error {
	customer := &clients.CustomerDTO{
		ID:    uint(customerID),
		Name:  fmt.Sprintf("Cliente %d", customerID),
		CPF:   uint(customerID * 1000),
		Email: fmt.Sprintf("cliente%d@example.com", customerID),
	}
	tc.customers[uint(customerID)] = customer

	tc.customerClient.EXPECT().
		ValidateCustomer(mock.Anything, uint(customerID)).
		Return(true, nil).
		Maybe()

	tc.customerClient.EXPECT().
		GetCustomer(mock.Anything, uint(customerID)).
		Return(customer, nil).
		Maybe()

	return nil
}

func (tc *testContext) oClienteComIDExisteComNome(customerID int, name string) error {
	customer := &clients.CustomerDTO{
		ID:    uint(customerID),
		Name:  name,
		CPF:   uint(customerID * 1000),
		Email: fmt.Sprintf("cliente%d@example.com", customerID),
	}
	tc.customers[uint(customerID)] = customer

	tc.customerClient.EXPECT().
		ValidateCustomer(mock.Anything, uint(customerID)).
		Return(true, nil).
		Maybe()

	tc.customerClient.EXPECT().
		GetCustomer(mock.Anything, uint(customerID)).
		Return(customer, nil).
		Maybe()

	return nil
}

func (tc *testContext) oClienteComIDNaoExiste(customerID int) error {
	tc.customerClient.EXPECT().
		ValidateCustomer(mock.Anything, uint(customerID)).
		Return(false, nil).
		Maybe()

	return nil
}

func (tc *testContext) oProdutoComIDExisteComPreco(productID int, price float64) error {
	product := &clients.ProductDTO{
		ID:          uint(productID),
		Name:        fmt.Sprintf("Produto %d", productID),
		Description: fmt.Sprintf("Descrição do produto %d", productID),
		Price:       float32(price),
		Category:    1,
		ImageLink:   "http://example.com/image.jpg",
	}
	tc.products[uint(productID)] = product

	// Mock será configurado dinamicamente quando chamarmos GetProducts
	tc.productClient.On("GetProducts", mock.Anything, mock.Anything).
		Return([]*clients.ProductDTO{product}, nil).
		Maybe()

	return nil
}

func (tc *testContext) oProdutoComIDExisteComNome(productID int, name string) error {
	product := &clients.ProductDTO{
		ID:          uint(productID),
		Name:        name,
		Description: fmt.Sprintf("Descrição de %s", name),
		Price:       10.0,
		Category:    1,
		ImageLink:   "http://example.com/image.jpg",
	}
	tc.products[uint(productID)] = product

	// Mock será configurado dinamicamente quando chamarmos GetProducts
	tc.productClient.On("GetProducts", mock.Anything, mock.Anything).
		Return([]*clients.ProductDTO{product}, nil).
		Maybe()

	return nil
}

func (tc *testContext) oProdutoComIDNaoExiste(productID int) error {
	// Produto não está no map, GetProducts retornará lista vazia para esse ID
	return nil
}

func (tc *testContext) euCriarUmPedidoParaOClienteComOsSeguintesProdutos(customerID int, table *godog.Table) error {
	tc.orderProducts = []*dto.AddOrderProductDto{}
	tc.totalAmount = 0

	for i, row := range table.Rows {
		if i == 0 {
			continue // Skip header
		}

		var prodID uint
		var qty uint
		var price float32
		fmt.Sscanf(row.Cells[0].Value, "%d", &prodID)
		fmt.Sscanf(row.Cells[1].Value, "%d", &qty)
		fmt.Sscanf(row.Cells[2].Value, "%f", &price)

		tc.orderProducts = append(tc.orderProducts, &dto.AddOrderProductDto{
			ProductId: prodID,
			Quantity:  qty,
			Price:     price,
		})
		tc.totalAmount += price * float32(qty)
	}

	// Mock do repositório
	tc.orderRepo.EXPECT().
		AddOrder(mock.Anything).
		Return(&entities.OrderEntity{ID: 1, CustomerId: uint(customerID), TotalAmount: tc.totalAmount}, nil).
		Maybe()

	tc.orderProductRepo.EXPECT().
		AddOrderProduct(mock.Anything).
		Return(nil).
		Maybe()

	tc.orderStatusRepo.EXPECT().
		AddOrderStatus(mock.Anything).
		Return(nil).
		Maybe()

	command := commands.NewAddOrderCommand(uint(customerID), tc.totalAmount, tc.orderProducts)
	tc.resultOrderID, tc.resultError = tc.addOrderUseCase.Execute(command)

	return nil
}

func (tc *testContext) euCriarUmPedidoParaOClienteComProduto(customerID, productID int) error {
	tc.orderProducts = []*dto.AddOrderProductDto{
		{ProductId: uint(productID), Quantity: 1, Price: 10.0},
	}
	tc.totalAmount = 10.0

	tc.orderRepo.EXPECT().
		AddOrder(mock.Anything).
		Return(&entities.OrderEntity{ID: 1, CustomerId: uint(customerID), TotalAmount: tc.totalAmount}, nil).
		Maybe()

	tc.orderProductRepo.EXPECT().
		AddOrderProduct(mock.Anything).
		Return(nil).
		Maybe()

	tc.orderStatusRepo.EXPECT().
		AddOrderStatus(mock.Anything).
		Return(nil).
		Maybe()

	command := commands.NewAddOrderCommand(uint(customerID), tc.totalAmount, tc.orderProducts)
	tc.resultOrderID, tc.resultError = tc.addOrderUseCase.Execute(command)

	return nil
}

func (tc *testContext) euCriarUmPedidoParaOClienteComProdutoInexistente(customerID, productID int) error {
	return tc.euCriarUmPedidoParaOClienteComProduto(customerID, productID)
}

func (tc *testContext) oPedidoDeveSerCriadoComSucesso() error {
	if tc.resultError != nil {
		return fmt.Errorf("esperava sucesso mas recebeu erro: %v", tc.resultError)
	}
	if tc.resultOrderID == "" {
		return fmt.Errorf("esperava ID do pedido mas recebeu string vazia")
	}
	return nil
}

func (tc *testContext) oTotalDoPedidoDeveSer(expectedTotal float64) error {
	if tc.totalAmount != float32(expectedTotal) {
		return fmt.Errorf("total esperado %.2f mas foi %.2f", expectedTotal, tc.totalAmount)
	}
	return nil
}

func (tc *testContext) oStatusDoPedidoDeveSer(expectedStatus int) error {
	// Status inicial é sempre 1 (Recebido) ao criar
	if expectedStatus != 1 {
		return fmt.Errorf("status esperado %d mas o inicial é sempre 1", expectedStatus)
	}
	return nil
}

func (tc *testContext) aCriacaoDoPedidoDeveFalhar() error {
	if tc.resultError == nil {
		return fmt.Errorf("esperava erro mas a operação foi bem-sucedida")
	}
	return nil
}

func (tc *testContext) oErroDeveConter(expectedMsg string) error {
	if tc.resultError == nil {
		return fmt.Errorf("esperava erro com mensagem '%s' mas não houve erro", expectedMsg)
	}
	if !contains(tc.resultError.Error(), expectedMsg) {
		return fmt.Errorf("esperava erro contendo '%s' mas recebeu '%s'", expectedMsg, tc.resultError.Error())
	}
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	tc := &testContext{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		tc.reset()
		return ctx, nil
	})

	// Contexto
	ctx.Step(`^que o serviço de clientes está disponível$`, tc.oServicoDeClientesEstaDisponivel)
	ctx.Step(`^o serviço de produtos está disponível$`, tc.oServicoDeProdutosEstaDisponivel)

	// Clientes
	ctx.Step(`^que o cliente com ID (\d+) existe$`, tc.oClienteComIDExiste)
	ctx.Step(`^o cliente com ID (\d+) existe com nome "([^"]*)"$`, tc.oClienteComIDExisteComNome)
	ctx.Step(`^que o cliente com ID (\d+) não existe$`, tc.oClienteComIDNaoExiste)

	// Produtos
	ctx.Step(`^o produto com ID (\d+) existe com preço ([\d.]+)$`, tc.oProdutoComIDExisteComPreco)
	ctx.Step(`^o produto com ID (\d+) existe com nome "([^"]*)"$`, tc.oProdutoComIDExisteComNome)
	ctx.Step(`^o produto com ID (\d+) não existe$`, tc.oProdutoComIDNaoExiste)

	// Ações
	ctx.Step(`^eu criar um pedido para o cliente (\d+) com os seguintes produtos:$`, tc.euCriarUmPedidoParaOClienteComOsSeguintesProdutos)
	ctx.Step(`^eu criar um pedido para o cliente (\d+) com produto (\d+)$`, tc.euCriarUmPedidoParaOClienteComProduto)
	ctx.Step(`^eu criar um pedido para o cliente (\d+) com produto inexistente (\d+)$`, tc.euCriarUmPedidoParaOClienteComProdutoInexistente)

	// Asserções
	ctx.Step(`^o pedido deve ser criado com sucesso$`, tc.oPedidoDeveSerCriadoComSucesso)
	ctx.Step(`^o total do pedido deve ser ([\d.]+)$`, tc.oTotalDoPedidoDeveSer)
	ctx.Step(`^o status do pedido deve ser (\d+)$`, tc.oStatusDoPedidoDeveSer)
	ctx.Step(`^a criação do pedido deve falhar$`, tc.aCriacaoDoPedidoDeveFalhar)
	ctx.Step(`^o erro deve conter "([^"]*)"$`, tc.oErroDeveConter)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
