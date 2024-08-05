# Usa uma imagem base com Go para build
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

# Copia o binário compilado e certificados de CA
COPY --from=builder /go-api /go-api
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Define a porta que a aplicação vai expor
EXPOSE 8080

# Define o comando a ser executado
ENTRYPOINT ["/go-api"]
