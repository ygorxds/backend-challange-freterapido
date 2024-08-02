# Usa uma imagem base com Go
FROM golang:1.20 AS builder

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-api main.go

# Usa uma imagem base mínima para rodar o binário
FROM scratch

# Copia o binário compilado
COPY --from=builder /go-api /go-api

# Define o comando a ser executado
ENTRYPOINT ["/go-api"]
