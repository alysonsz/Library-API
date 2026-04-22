# LIBRARY-API

## Descrição
O **LIBRARY-API** é um sistema backend desenvolvido em Go para gerenciar o catálogo de uma livraria. O sistema permite adicionar, atualizar, listar e remover livros, autores e empréstimos. Ele suporta **SQLite** e **PostgreSQL**, oferece uma API RESTful completa, CLI para comandos e está containerizado com Docker.

## Funcionalidades

- **Livros**: CRUD completo, busca paginada com filtros (gênero, autor, ano)
- **Autores**: CRUD completo, listar livros por autor
- **Empréstimos**: Registrar (com verificação de disponibilidade), retornar e listar atrasados
- **CLI**: Comandos `search`, `simulate`, `authors`, `loans`
- **Logging**: Structured JSON logging com `slog`
- **Graceful shutdown**: Finalização limpa do servidor HTTP
- **Health check**: Endpoint `/health` para monitoramento
- **CORS**: Suporte para requisições cross-origin
- **Validação**: Validação de entrada para livros, autores e empréstimos

## Requisitos

- Go instalado localmente
- SQLite3 (padrão) ou Docker com PostgreSQL

### Dependências

- `github.com/mattn/go-sqlite3`
- `github.com/jackc/pgx/v5`

## Como Executar

### SQLite (padrão)

```bash
go mod tidy
go run ./cmd/gobook
```

### PostgreSQL via Docker

```bash
docker-compose up --build
```

### Variáveis de Ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `PORT` | `8080` | Porta do servidor |
| `DB_DRIVER` | `sqlite3` | `sqlite3` ou `postgres` |
| `DB_PATH` | `./books.db` | Caminho do SQLite |
| `DB_URL` | (vazio) | Connection string do PostgreSQL |
| `LOG_LEVEL` | `info` | `info` ou `debug` |

### Exemplo com PostgreSQL

```bash
export DB_DRIVER=postgres
export DB_URL=postgres://user:pass@localhost:5432/library?sslmode=disable
go run ./cmd/gobook
```

## CLI

```bash
go run ./cmd/gobook search "Go"
go run ./cmd/gobook simulate 1 2 3
go run ./cmd/gobook authors list
go run ./cmd/gobook authors create "George Orwell" "Bio" 1903
go run ./cmd/gobook authors delete 1
go run ./cmd/gobook loans list
go run ./cmd/gobook loans overdue
go run ./cmd/gobook loans return 1
```

## Endpoints

### Livros

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/books?page=1&pageSize=10&genre=Fiction&authorId=1&publicationYear=2020` | Listar com filtros e paginação |
| POST | `/books` | Criar livro |
| GET | `/books/{id}` | Buscar por ID |
| PUT | `/books/{id}` | Atualizar |
| DELETE | `/books/{id}` | Remover |

### Autores

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/authors` | Listar |
| POST | `/authors` | Criar |
| GET | `/authors/{id}` | Buscar por ID |
| PUT | `/authors/{id}` | Atualizar |
| DELETE | `/authors/{id}` | Remover |
| GET | `/authors/{id}/books` | Livros do autor |

### Empréstimos

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/loans` | Listar |
| POST | `/loans` | Criar empréstimo |
| GET | `/loans/{id}` | Buscar por ID |
| POST | `/loans/{id}/return` | Devolver livro |
| GET | `/loans/overdue` | Empréstimos atrasados |

### Health

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/health` | Status do sistema e conectividade com banco |

### Swagger

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/swagger/index.html` | Interface interativa da documentação da API |

## Testes

```bash
go test ./... -v
```

Use o arquivo `test.http` para testar os endpoints da API.

## Swagger / OpenAPI

Para gerar a documentação Swagger, instale o `swag` CLI e rode:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/gobook/main.go
```

Após gerar, acesse `http://localhost:8080/swagger/index.html` para explorar a API interativamente.

Você também pode usar o Makefile:

```bash
make swag
```
