# API de Cotações de Frete

Bem-vindo ao repositório do projeto API de Cotações de Frete! Esta API RESTful, desenvolvida em Go, permite gerenciar cotações de frete e interagir com um banco de dados PostgreSQL.

## Descrição

A API de Cotações de Frete oferece os seguintes recursos:

- **Criação e Consulta de Cotações**: Envie novas cotações de frete e consulte cotações armazenadas.
- **Conexão com Banco de Dados PostgreSQL**: Utiliza um banco de dados PostgreSQL para armazenamento e recuperação de informações.
- **Testes Automatizados**: Inclui uma suíte de testes para assegurar a integridade e o correto funcionamento da API.

## Tecnologias Utilizadas

- **Go**: Linguagem de programação utilizada no desenvolvimento.
- **PostgreSQL**: Sistema de gerenciamento de banco de dados.
- **sqlx**: Biblioteca para facilitar interações com o PostgreSQL.
- **godotenv**: Carrega variáveis de ambiente de arquivos `.env`.
- **testing**: Biblioteca para desenvolvimento orientado a testes (TDD).

## Estrutura do Projeto

- **`main.go`**: Ponto de entrada da aplicação, contendo os modelos, roteamento, e controladores.
- **`main_test.go`**: Contém testes automatizados para validar o funcionamento da API.

## Configuração

Para configurar e iniciar o projeto, execute o seguinte comando:

```bash ```
docker-compose up --build


## Teste

Para testar, execute o seguinte comando:

```bash```
go test
