# Rate Limiter em Go

Este projeto implementa um rate limiter em Go que pode ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

## Configuração

### Variáveis de Ambiente

Defina as variáveis de ambiente no arquivo `.env`:

RATE_LIMIT_IP=10
RATE_LIMIT_TOKEN=100
BLOCK_DURATION=300
REDIS_ADDR=redis:6379
REDIS_PASSWORD=


### Docker

Use Docker Compose para subir o Redis e a aplicação:

```sh
docker-compose up --build -d
```
A aplicação estará disponível em http://localhost:8080.

### Uso
A limitação de requisições é aplicada com base no endereço IP ou no token de acesso passado no cabeçalho API_KEY. Se o limite for excedido, a resposta será HTTP 429 - Too Many Requests.

### Testes
Execute os testes com o comando:

```sh
docker-compose run --rm --no-deps --build --no-cache app go test ./tests
```

### Estrutura do Projeto
`cmd/main.go`: Configuração do servidor web.

`internal/middleware/middleware.go`: Middleware de rate limiting.

`internal/ratelimiter/ratelimiter.go`: Lógica do rate limiter.

`internal/redisclient/redisclient.go`: Inicialização e configuração do cliente Redis.

`Dockerfile`: Dockerfile para a aplicação.

`docker-compose.yml`: Configuração do Docker Compose para Redis e a aplicação.

`tests/middleware_test.go`: Testes automatizados.
