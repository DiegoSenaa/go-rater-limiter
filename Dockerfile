# Utilizando uma imagem oficial do Go
FROM golang:alpine3.20

# Criando um diretório de trabalho
WORKDIR /app

# Copiando os arquivos go.mod e go.sum e baixando as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiando o restante do código-fonte da aplicação
COPY . .

# Construindo o executável da aplicação
RUN go build -o main ./cmd/main.go

# Expondo a porta 8080
EXPOSE 8080

# Comando para iniciar a aplicação
CMD ["./main"]
