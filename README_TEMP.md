
# LIBRARY-API

## Descrição
O **LIBRARY-API** é um sistema backend desenvolvido em Go para gerenciar o catálogo de uma livraria. O sistema permite adicionar, atualizar, listar e remover livros do catálogo, com suporte para diferentes informações, como título, autor e gênero. Ele usa uma base de dados SQLite para armazenar os dados e oferece uma API RESTful para interagir com o sistema.

## Funcionalidades

- **Listar Livros**: Exibe todos os livros registrados no sistema.
- **Adicionar Livros**: Permite o cadastro de novos livros no catálogo.
- **Atributos do livro**: título, autor, gênero, quantidade de páginas, ano de publicação.
- **Atualizar Livro**: Permite a modificação dos detalhes de um livro existente.
- **Remover Livro**: Permite excluir um livro específico do sistema.

## Requisitos

- Go 1.22.5 ou superior
- SQLite3 para o banco de dados

### Dependências

- `github.com/mattn/go-sqlite3 v1.14.23`

## Como Executar

1. Clone o Repositório:
    ``` bash
    git clone https://github.com/alysonsz/library-api.git
    cd library-api
    ```
2. Instale as dependências:
   ```bash
   go mod tidy
   ```
3. Configurar o Banco de Dados:
O projeto utiliza o SQLite, então certifique-se de que você tenha o arquivo books.database no diretório correto. Caso ainda não exista, crie-o:
    ```bash
    touch books.database
    ```   
4. Inicie o servidor:
   ```bash
   go run main.go
   ```
5. O sistema estará rodando na porta `8080`.

## Endpoints

### Listar Livros
```http
GET /books
```
Retorna a lista de todos os livros cadastrados.

### Adicionar Livro
```http
POST /books
```
Corpo da requisição:
```json
{
  "title": "Nome do Livro",
  "author": "Autor do Livro",
  "genre": "Gênero",
  "pages": "Número de páginas",
  "publicationYear": "Ano de publicação"
}
```

### Atualizar Livro
```http
PUT /books/{id}
```
Parâmetro: `id` do livro a ser atualizado.  
Corpo da requisição:
```json
{
  "title": "Novo Nome",
  "author": "Novo Autor",
  "genre": "Novo Gênero",
  "pages": "Novo número de páginas",
  "publicationYear": "Novo ano de publicação"
}
```

### Remover Livro
```http
DELETE /books/{id}
```
Parâmetro: `id` do livro a ser removido.

## Testes

Use o arquivo `test.http` para testar os endpoints da API com exemplos de requisições HTTP.
