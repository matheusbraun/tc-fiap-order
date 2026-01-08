# Tech Challenge - Order Microservice

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=coverage)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=bugs)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order)

Um microsserviço standalone para gerenciamento de pedidos seguindo os princípios de Clean Architecture, extraído de uma aplicação monolítica.

## Índice

- [Sobre](#sobre)
- [Funcionalidades](#funcionalidades)
- [Tecnologias](#tecnologias)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Configuração](#configuração)
- [Uso](#uso)
- [Arquivos HTTP](#arquivos-http)
- [Testes com BDD](#testes-com-bdd)
- [Qualidade de Código](#qualidade-de-código)
- [Troubleshooting](#troubleshooting)

## Sobre

Este microserviço faz parte de uma arquitetura de microserviços para gestão de restaurantes. Ele foi desenvolvido em Golang com **PostgreSQL** como banco de dados, implementando **Clean Architecture** para separação clara entre regras de negócio e infraestrutura.

### Arquitetura Clean

O projeto segue os princípios da Clean Architecture, organizando o código em camadas bem definidas:

- **Domain (Domínio)**: Contém as entidades e regras de negócio centrais
- **Use Cases (Casos de Uso)**: Implementa a lógica de aplicação e orquestra as operações
- **Controllers**: Gerenciam o fluxo de dados entre a camada de apresentação e casos de uso
- **Infrastructure (Infraestrutura)**: Implementa detalhes técnicos como persistência PostgreSQL e APIs REST
- **Presenters**: Formatam os dados para apresentação com enriquecimento de dados externos

## Funcionalidades

### API REST de Gerenciamento de Pedidos

- ✅ **Criação de Pedidos**: Registre novos pedidos com produtos e valores
- ✅ **Consulta de Pedidos**: Busque pedidos individuais ou liste todos os pedidos ativos
- ✅ **Rastreamento de Status**: Acompanhe o status do pedido em tempo real
- ✅ **Atualização de Status**: Atualize o status do pedido através do ciclo de vida
- ✅ **Enriquecimento de Dados**: Integração com serviços de clientes e produtos
- ✅ **Degradação Graciosa**: Continua operando mesmo se serviços externos falharem
- ✅ **API RESTful**: Interface padronizada seguindo boas práticas REST
- ✅ **Documentação Swagger**: API totalmente documentada com OpenAPI 3.0

### Infraestrutura e DevOps

- ✅ **PostgreSQL 15**: Banco de dados relacional robusto
- ✅ **Docker & Docker Compose**: Containerização completa da aplicação
- ✅ **Kubernetes Ready**: Manifestos K8s para orquestração em clusters
- ✅ **Horizontal Pod Autoscaler**: Auto-scaling baseado em CPU e memória
- ✅ **Health Checks**: Endpoints de monitoramento de saúde da aplicação

### Qualidade e Arquitetura

- ✅ **Clean Architecture**: Separação clara de responsabilidades e camadas
- ✅ **Dependency Injection**: Gerenciamento com Uber FX
- ✅ **Testes Unitários**: 85 testes com cobertura ≥ 80%
- ✅ **BDD Pattern**: Testes com padrão Given/When/Then
- ✅ **SonarCloud**: Análise contínua de qualidade de código
- ✅ **Mocks Automatizados**: Geração de mocks para testes isolados

### Banco de Dados

- **Tabelas**: `order`, `order_product`, `order_status`
- **Isolamento**: Sem chaves estrangeiras para serviços externos
- **Histórico**: Status do pedido mantém histórico completo
- **Migrations**: Criação automática de schema

## Tecnologias

- **Go 1.24.2+** - Linguagem de programação principal
- **PostgreSQL 15** - Banco de dados relacional
- **GORM v1.26.1** - ORM para Go
- **Chi Router v5.2.1** - Router HTTP leve e performático
- **Uber FX v1.23.0** - Framework de injeção de dependências
- **Testify v1.11.1** - Framework de testes
- **Mockery v2.53.5** - Geração de mocks
- **Swagger/OpenAPI** - Documentação da API
- **Docker & Docker Compose** - Containerização
- **Kubernetes** - Orquestração de containers

## Estrutura do Projeto

O projeto segue os princípios da **Clean Architecture**, organizando o código em camadas bem definidas:

```
cmd/api/                                # Entrada da aplicação (main.go)
internal/
  app/                                  # Inicialização e injeção de dependências
  infrastructure/
    clients/                            # Clientes HTTP para serviços externos
      customer_client.go                # Cliente do serviço de clientes
      product_client.go                 # Cliente do serviço de produtos
  order/                                # Domínio de Pedidos
    controller/                         # Controllers (orquestração)
      order_controller.go
      order_controller_impl.go
      order_controller_test.go
    domain/
      entities/                         # Entidades do domínio
        order.go
        order_product.go
        order_status.go
      repositories/                     # Interfaces dos repositórios
        order_repository.go
        order_product_repository.go
        order_status_repository.go
    infrastructure/
      api/                              # HTTP/REST API
        controller/
          order_api_controller.go       # Handlers HTTP
        dto/                            # Data Transfer Objects
          add_order_dto.go
          get_order_response_dto.go
          get_orders_response_dto.go
          get_orderstatus_response_dto.go
          update_order_status_request_dto.go
      persistence/                      # Data persistence
        order_repository_impl.go
        order_repository_impl_test.go
        order_product_repository_impl.go
        order_product_repository_impl_test.go
        order_status_repository_impl.go
        order_status_repository_impl_test.go
    presenter/                          # Presentation layer
      order_presenter.go
      order_presenter_impl.go
      order_presenter_test.go
    usecase/                            # Business logic use cases
      addOrder/
        add_order_use_case.go
        add_order_use_case_impl.go
        add_order_use_case_test.go
      getOrder/
        get_order_use_case.go
        get_order_use_case_impl.go
        get_order_use_case_test.go
      getOrders/
        get_orders_use_case.go
        get_orders_use_case_impl.go
        get_orders_use_case_test.go
      getOrderStatus/
        get_order_status_use_case.go
        get_orders_status_use_case_impl.go
        get_order_status_use_case_test.go
      updateOrderStatus/
        update_order_status_use_case.go
        update_order_status_use_case._impl.go
        update_order_status_use_case_test.go
      commands/                         # Command pattern objects
        add_order_command.go
        get_order_command.go
        get_orders_command.go
        get_order_status_command.go
        update_order_status_command.go
  shared/                               # Shared utilities
    config/                             # Configuration management
    httpclient/                         # HTTP client with retry logic
pkg/                                    # Public shared packages
  rest/                                 # REST utilities
  storage/postgres/                     # PostgreSQL connection
mocks/                                  # Auto-generated mocks
  order/
    controller/
    domain/repositories/
    presenter/
    usecase/
  infrastructure/clients/
  shared/httpclient/
docs/                                   # Swagger documentation (generated)
http/                                   # HTTP test files (REST Client)
k8s/                                    # Kubernetes manifests
scripts/                                # Utility scripts
  coverage.sh
  coverage.ps1
.github/workflows/                      # CI/CD pipelines
go.mod                                  # Go module dependencies
go.sum                                  # Dependency checksums
Makefile                                # Build automation
.mockery.yaml                           # Mock generation config
docker-compose.yml                      # Local development setup
Dockerfile                              # Container image
README.md                               # Project documentation
```

## Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` baseado em `.env.example`:

```bash
# Configuração do Servidor
SERVER_PORT=8080

# Configuração do Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=order_user
DB_PASSWORD=order_pass
DB_NAME=order_db
DB_SSLMODE=disable

# URLs de Serviços Externos
CUSTOMER_SERVICE_URL=http://localhost:8081
PRODUCT_SERVICE_URL=http://localhost:8082

# Configuração do Cliente HTTP
HTTP_CLIENT_TIMEOUT_SECONDS=30
HTTP_CLIENT_RETRY_COUNT=3
HTTP_CLIENT_RETRY_BACKOFF_MS=100
```

### Desenvolvimento Local

#### Opção 1: Usando Make

```bash
# Instalar dependências
go mod download

# Executar a aplicação
make run

# Compilar a aplicação
make build

# Executar testes
make test

# Gerar cobertura de testes
make coverage

# Gerar mocks
make mocks
```

#### Opção 2: Comandos Go Diretos

```bash
# Executar a aplicação
go run ./cmd/api/main.go

# Compilar a aplicação
go build -o bin/order-service ./cmd/api

# Executar testes
go test ./... -v
```

### Desenvolvimento com Docker

```bash
# Iniciar todos os serviços (banco de dados + aplicação)
make docker-up

# Ou usando docker-compose diretamente
docker-compose up -d

# Parar serviços
make docker-down
```

### Deploy no Kubernetes

```bash
# Aplicar configurações
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/hpa.yaml

# Verificar status do deployment
kubectl get pods -l app=order-service
kubectl get svc order-service
```

## Uso

### Endpoints da API

#### 1. Criar Pedido
```bash
POST /v1/order
Content-Type: application/json

{
  "customerId": 1,
  "totalAmount": 150.00,
  "products": [
    {
      "productId": 10,
      "quantity": 2,
      "price": 50.00
    },
    {
      "productId": 20,
      "quantity": 1,
      "price": 50.00
    }
  ]
}
```

**Resposta (201 Created):**
```json
{
  "id": 123,
  "created_at": "2026-01-07T23:00:00Z",
  "total_amount": 150.00,
  "customer_id": 1,
  "customer": {
    "id": 1,
    "name": "João Silva",
    "email": "joao@example.com",
    "cpf": 12345678901
  },
  "products": [
    {
      "product_id": 10,
      "price": 50.00,
      "quantity": 2,
      "name": "Hambúrguer",
      "description": "Hambúrguer artesanal",
      "category": 1
    }
  ],
  "status": [
    {
      "id": 1,
      "current_status": 1,
      "current_status_description": "Recebido",
      "order_id": 123
    }
  ]
}
```

#### 2. Buscar Pedido por ID
```bash
GET /v1/order/123
```

**Resposta (200 OK):** (mesmo formato do POST)

#### 3. Listar Pedidos Ativos
```bash
GET /v1/order
```

**Resposta (200 OK):**
```json
{
  "orders": [
    {
      "id": 123,
      "created_at": "2026-01-07T23:00:00Z",
      "total_amount": 150.00,
      "customer_id": 1,
      "products": [...],
      "status": [...]
    }
  ]
}
```

#### 4. Buscar Status do Pedido
```bash
GET /v1/order/123/status
```

**Resposta (200 OK):**
```json
{
  "id": 1,
  "created_at": "2026-01-07T23:00:00Z",
  "current_status": 2,
  "current_status_description": "Em preparação",
  "order_id": 123
}
```

#### 5. Atualizar Status do Pedido
```bash
PUT /v1/order/123/status
Content-Type: application/json

{
  "status": 3
}
```

**Resposta (200 OK)**

### Ciclo de Vida do Status do Pedido

1. **Recebido (1)** - Pedido recebido
2. **Em preparação (2)** - Sendo preparado
3. **Pronto (3)** - Pronto para retirada
4. **Finalizado (4)** - Pedido concluído

### Documentação Swagger

Swagger UI disponível em: `http://localhost:8080/swagger/`

## Arquivos HTTP

No diretório [`http/`](http/) estão arquivos `.http` com exemplos prontos para testar a API usando a extensão [REST Client para VS Code](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).

**Como usar:**
1. Abra o arquivo `.http` no VS Code
2. Ajuste a variável `baseUrl` se necessário (ex: `@baseUrl = http://localhost:8080/`)
3. Clique em "Send Request" para executar e ver a resposta

Os exemplos cobrem todos os endpoints da API de pedidos.

## Testes com BDD

Este projeto implementa **Behavior-Driven Development (BDD)** usando [testify/suite](https://pkg.go.dev/github.com/stretchr/testify/suite).

- **107 testes** em 9 camadas (clients, repository, use case, controller, presenter)
- Padrão **Given/When/Then** para clareza e legibilidade
- Mocks gerados automaticamente com Mockery
- Testes descritivos que funcionam como documentação

### Executar Testes

```bash
# Executar todos os testes
make test

# Executar com cobertura
make coverage

# Gerar relatório HTML de cobertura
make coverage-report
```

### Exemplo de Teste BDD

```go
// Feature: Order Repository - Add Order
// Scenario: Create a new order successfully

func (suite *OrderRepositoryTestSuite) Test_AddOrder_WithValidData_ShouldCreateSuccessfully() {
    // GIVEN a valid order entity
    order := &entities.OrderEntity{
        CustomerId:  1,
        TotalAmount: 100.50,
    }

    // WHEN the order is added to the repository
    result, err := suite.repository.AddOrder(order)

    // THEN the operation should complete without errors
    assert.NoError(suite.T(), err)
    // AND the order should have an ID assigned
    assert.NotZero(suite.T(), result.ID)
}
```

## Qualidade de Código

Este projeto utiliza **SonarCloud** para análise contínua de qualidade de código, segurança e cobertura de testes.

### Métricas Monitoradas

| Métrica | Status | Objetivo |
|---------|--------|----------|
| **Quality Gate** | [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order) | ✅ Passed |
| **Coverage** | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=coverage)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order) | ≥ 80% |
| **Bugs** | [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=bugs)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order) | 0 |
| **Code Smells** | [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order) | < 10 |
| **Security** | [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=matheusbraun_tc-fiap-order&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=matheusbraun_tc-fiap-order) | A |

### Análise Automática

A análise de código é executada automaticamente via **GitHub Actions** em:
- ✅ Cada push para branches `main` e `develop`
- ✅ Todos os Pull Requests

### Gerar Coverage Localmente

**Windows (PowerShell):**
```powershell
.\scripts\coverage.ps1
```

**Linux/Mac (Bash):**
```bash
chmod +x scripts/coverage.sh
./scripts/coverage.sh
```

Isso gera:
- `coverage.out` - Formato para SonarCloud
- `coverage.html` - Visualização no browser

## Troubleshooting

### Problemas Comuns

**Problema**: Falha na conexão com o banco de dados
```bash
# Verificar se o banco de dados está rodando
docker ps | grep postgres

# Verificar variáveis de ambiente
env | grep DB_
```

**Problema**: Falhas em chamadas a serviços externos
```bash
# Verificar URLs dos serviços
curl $CUSTOMER_SERVICE_URL/health
curl $PRODUCT_SERVICE_URL/health

# Verificar conectividade de rede
```

**Problema**: Falha na compilação
```bash
# Limpar e recompilar
make clean
go mod tidy
go build ./...
```

## Integração com Serviços Externos

### Serviço de Clientes

- **URL Base**: Configurada via `CUSTOMER_SERVICE_URL`
- **Usado para**: Validação de clientes e enriquecimento de dados
- **Endpoints**:
  - `GET /v1/customer/{id}` - Buscar detalhes do cliente
- **Degradação Graciosa**: Se falhar, retorna pedido sem dados do cliente

### Serviço de Produtos

- **URL Base**: Configurada via `PRODUCT_SERVICE_URL`
- **Usado para**: Validação de produtos e enriquecimento de dados
- **Endpoints**:
  - `GET /v1/product/{id}` - Buscar detalhes do produto
- **Degradação Graciosa**: Se falhar, retorna pedido sem dados enriquecidos do produto

### Lógica de Retry

Clientes HTTP implementam retry com backoff exponencial:
- **Tentativas Padrão**: 3 tentativas
- **Backoff**: 100ms base (exponencial)
- **Timeout**: 30 segundos por requisição

## Performance

### Limites de Recursos (Kubernetes)

- **Requests**: 100m CPU, 128Mi Memória
- **Limits**: 200m CPU, 256Mi Memória

### Auto-scaling (HPA)

- **Réplicas Mínimas**: 2
- **Réplicas Máximas**: 10
- **CPU Alvo**: 70%
- **Memória Alvo**: 80%

## Migração do Monolito

### Principais Mudanças

1. **Removidas Chaves Estrangeiras**: Pedidos não possuem mais constraints FK para tabelas de cliente/produto
2. **Adicionados Clientes HTTP**: Dados de cliente e produto buscados via HTTP
3. **Banco de Dados Isolado**: Banco de dados PostgreSQL separado para o serviço de pedidos
4. **Presenters Atualizados**: Lógica de enriquecimento movida do preload GORM para chamadas HTTP
5. **Gerenciamento de Configuração**: Configuração centralizada com variáveis de ambiente

### Breaking Changes

- Entidade Customer não está mais embutida em OrderEntity
- Entidade Product não está mais embutida em OrderProductEntity
- Respostas podem ter dados de cliente/produto nulos se serviços externos falharem
