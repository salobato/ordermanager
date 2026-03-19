# 🧾 Order Manager

API para gerenciamento de pedidos, construída com foco em **Clean Architecture**, **boas práticas de design** e **testabilidade**.

---

## Executando o projeto

### Pré-requisitos

* Docker
* Docker Compose

---

### Subindo a aplicação

```bash
docker-compose up --build
```

A aplicação estará disponível em:

```
http://localhost:8080
```

Todas as dependências estão no container, não sendo necessária a instalação de nenhuma dependência como MongoDB ou RabbitMQ.

---

### Healthcheck

```http
GET /health
```

---

## Stack utilizada

### Backend

* Go 1.26.1 (Golang)
* Gin (framework HTTP)

### Banco de dados

* MongoDB

### Mensageria

* RabbitMQ

### Testes

* Testify (assert + mock)

### Infraestrutura

* Docker
* Docker Compose

---

## Funcionalidades implementadas

* Criar pedido
* Atualizar status do pedido
* Buscar pedido por ID
* Publicação de eventos via RabbitMQ
* Healthcheck com verificação de dependências

---

## Testes

O projeto possui cobertura de testes em múltiplos níveis:

### Testes unitários

* Use cases (regras de negócio)
* Entidades (validações e transições de status)

### Testes com mocks

* Repositórios
* Publisher (RabbitMQ)

### Testes de integração

* MongoDB (tabela de testes controlados que é gerada e apagada toda vez que a suíte de testes é executada)
* RabbitMQ (publicação e consumo de mensagens)

## Como executar a suíte de testes?

Para visualizar o coverage:
```bash
go test -cover ./... 
```

Ou com a flag -v para mostrar cada teste:
```bash
go test -v -cover ./... 
```

O output mais recente:
```bash
go test -cover ./...
        github.com/salobato/ordermanager/cmd/api                coverage: 0.0% of statements
        github.com/salobato/ordermanager/internal/adapter/api/gin/health                coverage: 0.0% of statements
ok      github.com/salobato/ordermanager/internal/adapter/api/gin/order 0.508s  coverage: 87.9% of statements
ok      github.com/salobato/ordermanager/internal/adapter/messaging/rabbitmq    1.660s  coverage: 77.8% of statements
ok      github.com/salobato/ordermanager/internal/adapter/repository/mongo      1.536s  coverage: 71.4% of statements
        github.com/salobato/ordermanager/internal/adapter/repository/mongo/models               coverage: 0.0% of statements
ok      github.com/salobato/ordermanager/internal/core/entity   0.792s  coverage: 97.1% of statements
?       github.com/salobato/ordermanager/internal/core/publisher        [no test files]
?       github.com/salobato/ordermanager/internal/core/repository       [no test files]
ok      github.com/salobato/ordermanager/internal/core/usecase  1.059s  coverage: 96.0% of statements
        github.com/salobato/ordermanager/pkg/config             coverage: 0.0% of statements
```
---

## Principais decisões técnicas

### 1. Separação por camadas (Clean Architecture + Hexagonal Architecture)

O projeto foi estruturado em camadas bem definidas:

```
internal/
  core/
    entity/
    usecase/
    repository/
  adapter/
    database/
    messaging/
    http/
```

**Motivação:**

* Baixo acoplamento
* Alta testabilidade
* Independência de frameworks

---

### 2. Uso de DTOs no layer HTTP

Os handlers não retornam diretamente as entidades do domínio.

Foi criada uma camada de **DTO (Data Transfer Object)** para:

* Padronizar respostas (snake_case)
* Evitar vazamento de detalhes internos (ex: ObjectID do Mongo)
* Permitir evolução da API sem impactar o domínio

---

### 3. Mensageria desacoplada

O envio de eventos (RabbitMQ) é feito via interface:

```go
type EventPublisher interface {
    PublishOrderStatusChanged(ctx context.Context, event OrderEvent) error
}
```

**Benefícios:**

* Facilita testes com mocks
* Permite troca de broker sem impactar o core

---

### 4. Geração de OrderNumber via contador

Foi implementado um contador no MongoDB para gerar números sequenciais:

```
ORD-2026-000001
```

Utilizando:

* `findOneAndUpdate` com `$inc`
* Upsert automático

---

### 5. Regras de negócio na entidade

A lógica de transição de status está encapsulada na entidade:

* Não é permitido pular etapas
* Pedido entregue não pode ser alterado

Isso garante **consistência do domínio independentemente da origem da chamada**.

---

### 6. Testabilidade como prioridade

* Interfaces para dependências externas
* Mocks para isolamento
* Testes de integração com containers reais

---

### 7. Ambiente isolado com Docker 🐳

Toda a stack roda via `docker-compose`, incluindo:

* API
* MongoDB
* RabbitMQ

---

## Exemplos de requisição

### Criar pedido

```http
POST /orders
```

```json
{
  "customer_id": "65f0c3a8b9d4e2a1f0c9b123",
  "total": 200.5
}
```

---

### Atualizar status

```http
PATCH /orders/:id/status
```

```json
{
  "status": "em_processamento"
}
```

---

### Buscar por ID

```http
GET /orders/:id
```

```json
{
    "id": "69bb46cd45a73e6cc3ec441c",
    "order_number": "ORD-2026-000003",
    "customer_id": "65f0c3a8b9d4e2a1f0c9b123",
    "total": 110,
    "status": "em_processamento",
    "placed_at": "2026-03-19T00:43:57Z",
    "updated_at": "2026-03-19T00:49:06Z"
}

---
