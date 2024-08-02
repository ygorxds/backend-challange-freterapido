# Projeto API de Cotações de Frete

Bem-vindo ao repositório do projeto API de Cotações de Frete! Este projeto é uma API RESTful desenvolvida em Go que permite gerenciar cotações de frete e interagir com uma base de dados PostgreSQL.

## Descrição

Esta API oferece endpoints para criar e consultar cotações de frete. O projeto inclui:

- **Criação e Consulta de Cotações**: Permite o envio de cotações de frete e a consulta de cotações armazenadas.
- **Conexão com Banco de Dados PostgreSQL**: Utiliza um banco de dados PostgreSQL para armazenar e recuperar informações.
- **Testes Automatizados**: Inclui testes automatizados para garantir a integridade e o funcionamento da API.

## Tecnologias Utilizadas

- **Go**: Linguagem de programação principal do projeto.
- **PostgreSQL**: Sistema de gerenciamento de banco de dados utilizado.
- **sqlx**: Biblioteca para interações com o banco de dados PostgreSQL.
- **godotenv**: Biblioteca para carregar variáveis de ambiente de um arquivo `.env`.
- **testing **: Biblioteca para realizar testes (TDD).
- 
## Estrutura do Projeto

- **`main.go`**: Ponto de entrada principal da aplicação. Configura o roteamento e os handlers da API.
- **`handlers.go`**: Contém a lógica dos handlers para os endpoints da API.
- **`models.go`**: Define os modelos de dados utilizados na aplicação.
- **`main_test.go`**: Contém os testes automatizados para verificar o funcionamento da API.

## Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis de ambiente:

```env
DATABASE_URL=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=



