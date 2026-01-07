# Microsserviço de Pedidos

Um microsserviço standalone para gerenciamento de pedidos seguindo os princípios de Clean Architecture, extraído de uma aplicação monolítica.

## Visão Geral

Este microsserviço gerencia todas as operações relacionadas a pedidos, incluindo:
- Criação de pedidos com validação de produtos
- Recuperação de pedidos com dados enriquecidos de clientes e produtos
- Listagem de pedidos (excluindo pedidos finalizados)
- Rastreamento e atualização de status de pedidos

## Arquitetura

### Camadas da Clean Architecture

```
├── Camada de Domínio (Lógica de Negócio)
│   ├── Entidades: OrderEntity, OrderProductEntity, OrderStatusEntity
│   └── Interfaces de Repositórios
│
├── Camada de Casos de Uso (Lógica de Aplicação)
│   ├── AddOrder: Criar pedidos com validação
│   ├── GetOrder: Buscar pedido com dados enriquecidos
│   ├── GetOrders: Listar pedidos ativos
│   ├── GetOrderStatus: Obter status atual do pedido
│   └── UpdateOrderStatus: Atualizar status do pedido
│
├── Camada de Controlador (Orquestração de Fluxo)
│   └── OrderController: Coordena DTOs e casos de uso
│
├── Camada de Infraestrutura
│   ├── Persistência: Repositórios PostgreSQL (GORM)
│   ├── API: Handlers HTTP (Chi router)
│   ├── Clients: Clientes HTTP para serviços externos
│   └── DTOs: Objetos de Requisição/Resposta
│
└── Camada de Apresentação (Formatação de Dados)
    └── OrderPresenter: Converte entidades para DTOs
```

### Características do Microsserviço

- **Banco de Dados Isolado**: Banco de dados PostgreSQL dedicado (sem tabelas compartilhadas)
- **Sem Chaves Estrangeiras**: Referências aos serviços de cliente/produto apenas via IDs
- **Comunicação HTTP Síncrona**: Chamadas RESTful para serviços externos
- **Degradação Graciosa**: Continua operando se serviços externos estiverem indisponíveis
- **Clean Architecture**: Mantém separação de responsabilidades e testabilidade

## Stack Tecnológica

- **Linguagem**: Go 1.24.2+
- **Roteador HTTP**: Chi v5.2.1
- **ORM**: GORM v1.26.1
- **Banco de Dados**: PostgreSQL 15
- **Injeção de Dependência**: Uber FX v1.23.0
- **Testes**: Testify + Godog (BDD)
- **Documentação da API**: Swagger/OpenAPI

## Começando

### Pré-requisitos

- Go 1.24.2 ou superior
- PostgreSQL 15
- Docker e Docker Compose (opcional)
- Make (opcional, para comandos de conveniência)

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
make test-coverage

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

## Endpoints da API

### Pedidos

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/v1/order` | Criar um novo pedido |
| GET | `/v1/order/{orderId}` | Obter pedido por ID |
| GET | `/v1/order` | Listar todos os pedidos ativos |
| GET | `/v1/order/{orderId}/status` | Obter status do pedido |
| PUT | `/v1/order/{orderId}/status` | Atualizar status do pedido |

### Documentação da API

Swagger UI disponível em: `http://localhost:8080/swagger/`

## Ciclo de Vida do Status do Pedido

1. **Recebido (1)** - Pedido recebido
2. **Em preparação (2)** - Sendo preparado
3. **Pronto (3)** - Pronto
4. **Finalizado (4)** - Concluído

## Integração com Serviços Externos

### Serviço de Clientes

- **URL Base**: Configurada via `CUSTOMER_SERVICE_URL`
- **Usado para**: Validação de clientes e enriquecimento de dados
- **Endpoints**:
  - `GET /v1/customer/{id}` - Buscar detalhes do cliente

### Serviço de Produtos

- **URL Base**: Configurada via `PRODUCT_SERVICE_URL`
- **Usado para**: Validação de produtos e enriquecimento de dados
- **Endpoints**:
  - `GET /v1/product/{id}` - Buscar detalhes do produto
  - Busca em lote implementada na camada de cliente

## Schema do Banco de Dados

### Tabela Order
```sql
CREATE TABLE "order" (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amount FLOAT DEFAULT 0,
    customer_id INTEGER NOT NULL  -- Sem constraint FK
);
CREATE INDEX idx_order_customer_id ON "order"(customer_id);
```

### Tabela Order Product
```sql
CREATE TABLE order_product (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES "order"(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL,  -- Sem constraint FK
    price FLOAT NOT NULL,
    quantity INTEGER NOT NULL
);
CREATE INDEX idx_order_product_order_id ON order_product(order_id);
CREATE INDEX idx_order_product_product_id ON order_product(product_id);
```

### Tabela Order Status
```sql
CREATE TABLE order_status (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    current_status INTEGER NOT NULL,
    order_id INTEGER NOT NULL REFERENCES "order"(id) ON DELETE CASCADE
);
CREATE INDEX idx_order_status_order_id ON order_status(order_id);
```

## Testes

### Estrutura de Testes

```
tests/
├── features/           # Arquivos de feature Gherkin BDD
│   ├── add_order.feature
│   ├── get_order.feature
│   └── ...
├── steps/              # Definições de steps BDD
│   └── order_steps.go
├── mocks/              # Mocks gerados (mockery)
│   ├── OrderRepository.go
│   ├── CustomerClient.go
│   └── ProductClient.go
└── unit/               # Testes unitários
    ├── controller/
    ├── presenter/
    └── clients/
```

### Executando Testes

```bash
# Executar todos os testes
make test

# Executar apenas testes unitários
make test-unit

# Executar apenas testes BDD
make test-bdd

# Gerar relatório de cobertura
make test-coverage
```

### Meta de Cobertura de Testes

Meta: **80%** de cobertura em todas as camadas

## Tratamento de Erros

### Degradação Graciosa

O serviço implementa degradação graciosa para falhas de serviços externos:

- **Serviço de Clientes Fora**: Pedidos são retornados sem dados do cliente
- **Serviço de Produtos Fora**: Pedidos são retornados apenas com IDs de produtos (sem dados enriquecidos)
- **Falhas de Validação**: Pedidos não podem ser criados se os serviços estiverem indisponíveis durante a criação

### Lógica de Retry

Clientes HTTP implementam retry com backoff exponencial:
- **Tentativas Padrão**: 3 tentativas
- **Backoff**: 100ms base (exponencial)
- **Timeout**: 30 segundos por requisição

## Desenvolvimento

### Estrutura do Projeto

```
.
├── cmd/api/                    # Ponto de entrada da aplicação
├── internal/
│   ├── domain/                 # Entidades de negócio e interfaces
│   ├── usecase/                # Lógica de aplicação
│   ├── controller/             # Orquestração de fluxo
│   ├── infrastructure/         # Integrações externas
│   │   ├── api/                # Handlers HTTP
│   │   ├── clients/            # Clientes de serviços externos
│   │   └── persistence/        # Repositórios de banco de dados
│   ├── presenter/              # Formatação de dados
│   └── shared/                 # Utilitários compartilhados
│       ├── config/             # Gerenciamento de configuração
│       └── httpclient/         # Wrapper de cliente HTTP
├── pkg/
│   ├── rest/                   # Utilitários REST
│   └── storage/postgres/       # Conexão com banco de dados
├── tests/                      # Arquivos de teste
├── k8s/                        # Manifests Kubernetes
├── Dockerfile                  # Configuração Docker
├── docker-compose.yml          # Configuração Docker Compose
└── docs/                       # Documentação da API
```

### Adicionando Novas Funcionalidades

1. **Adicionar Entidade de Domínio/Interface de Repositório** (se necessário)
2. **Implementar Caso de Uso** com lógica de negócio
3. **Criar Método no Controlador** para orquestração
4. **Adicionar Handler da API** na camada de infraestrutura
5. **Atualizar Presenter** para formatação de resposta
6. **Escrever Testes** (BDD + Unitários)
7. **Atualizar Documentação Swagger**

### Geração de Código

```bash
# Gerar mocks
make mocks

# Gerar documentação Swagger (se necessário)
swag init -g cmd/api/main.go
```

## Monitoramento e Health Checks

### Liveness Probe
- **Endpoint**: `GET /swagger/doc.json`
- **Delay Inicial**: 30s
- **Período**: 10s

### Readiness Probe
- **Endpoint**: `GET /swagger/doc.json`
- **Delay Inicial**: 10s
- **Período**: 5s

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

## Contribuindo

1. Siga os princípios da Clean Architecture
2. Escreva testes para novas funcionalidades (BDD + Unitários)
3. Mantenha cobertura de testes acima de 80%
4. Atualize a documentação Swagger
5. Use commits convencionais

## Licença

[Sua Licença Aqui]

## Suporte

Para problemas e dúvidas, consulte o rastreador de issues.
