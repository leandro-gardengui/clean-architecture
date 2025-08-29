# Clean Architecture Orders API

Sistema de gerenciamento de pedidos implementado seguindo os princípios de Clean Architecture com múltiplas interfaces: REST API, gRPC e GraphQL.

## Arquitetura

Projeto estruturado em camadas seguindo Clean Architecture:

```
├── cmd/server/           # Ponto de entrada da aplicação
├── internal/
│   ├── domain/           # Camada de domínio (entities, interfaces)
│   ├── usecase/          # Casos de uso da aplicação
│   └── infra/            # Camada de infraestrutura
│       ├── database/     # Implementação do repositório
│       ├── web/          # Handlers REST
│       ├── grpc/         # Serviços gRPC
│       └── graphql/      # Resolvers GraphQL
├── migrations/           # Migrações do banco de dados
└── api/                  # Definições de APIs
```

## Funcionalidades

- ✅ Criar pedidos (CreateOrder)
- ✅ Listar pedidos (ListOrders)
- ✅ REST API (POST/GET /order)
- 🚧 gRPC Service (temporariamente desabilitado - requer protoc)
- ✅ GraphQL Query/Mutation
- ✅ Banco PostgreSQL com migrações
- ✅ Docker/Docker Compose

## Portas dos Serviços

| Serviço  | Porta  | Endpoint/URL               |
|----------|--------|----------------------------|
| REST API | 8000   | http://localhost:8000      |
| gRPC     | 50051  | localhost:50051            |
| GraphQL  | 8080   | http://localhost:8080      |
| PostgreSQL| 5432  | localhost:5432             |

## Como Executar

### 1. Executar com Docker Compose (Recomendado)

```bash
# Clone o repositório e navegue até o diretório
cd clean-architecture

# Execute com docker compose
docker compose up --build

# Para parar os serviços
docker compose down
```

### 2. Executar Localmente

```bash
# Instalar dependências
go mod tidy

# Executar PostgreSQL
docker run -d \
  --name postgres \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=orders \
  -p 5432:5432 \
  postgres:15

# Executar aplicação
go run cmd/server/main.go
```

## Testando as APIs

### REST API

```bash
# Criar pedido
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"price": 100.0, "tax": 10.0}'

# Listar pedidos
curl http://localhost:8000/order
```

### GraphQL

Acesse http://localhost:8080 para o playground GraphQL ou use:

```bash
# Criar pedido
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { createOrder(price: 150.0, tax: 15.0) { id price tax final_price } }"}'

# Listar pedidos
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ listOrders { id price tax final_price } }"}'
```

### gRPC

Use um cliente gRPC como evans ou grpcurl:

```bash
# Instalar grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Listar serviços
grpcurl -plaintext localhost:50051 list

# Criar pedido
grpcurl -plaintext -d '{"price": 200.0, "tax": 20.0}' \
  localhost:50051 pb.OrderService/CreateOrder

# Listar pedidos
grpcurl -plaintext -d '{}' \
  localhost:50051 pb.OrderService/ListOrders
```

## Arquivo de Teste

Use o arquivo `api.http` incluído no projeto para testar facilmente com VS Code REST Client ou similar.

## Tecnologias Utilizadas

- **Go 1.21**
- **PostgreSQL** - Banco de dados
- **Gin** - Framework web para REST API
- **gRPC** - Comunicação de alta performance (requer protoc)
- **GraphQL** - API flexível para consultas
- **Docker/Docker Compose** - Containerização
- **golang-migrate** - Migrações de banco

## Nota sobre gRPC

O serviço gRPC está temporariamente desabilitado devido a problemas com arquivos protobuf gerados manualmente. Para habilitar:

1. Instale o protoc: `sudo apt install protobuf-compiler` (Ubuntu/Debian)
2. Instale os plugins Go: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
3. `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
4. Execute: `bash scripts/generate_proto.sh`
5. Descomente as linhas gRPC em `cmd/server/main.go`