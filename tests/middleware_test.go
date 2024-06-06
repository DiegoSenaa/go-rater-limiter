package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DiegoSenaa/go-rater-limiter/internal/middleware"
	"github.com/DiegoSenaa/go-rater-limiter/internal/ratelimiter"
	"github.com/DiegoSenaa/go-rater-limiter/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Configurando variáveis de ambiente para os testes
	os.Setenv("RATE_LIMIT_IP", "5")
	os.Setenv("RATE_LIMIT_TOKEN", "10")
	os.Setenv("BLOCK_DURATION", "60") // 1 minuto para facilitar os testes
	os.Setenv("REDIS_ADDR", "redis:6379")
	os.Setenv("REDIS_PASSWORD", "")

	// Executando os testes
	code := m.Run()

	// Saindo com o código de status adequado
	os.Exit(code)
}

func setupRouter(rl *ratelimiter.RateLimiter) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RateLimitMiddleware(rl))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})
	return r
}

func clearStorage(s storage.Storage) {
	fmt.Println("Limpando o Storage")
	err := s.Clear()
	if err != nil {
		fmt.Println("Error clearing storage:", err)
	}
}

func TestRateLimitByIP(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	rateLimiter := ratelimiter.NewRateLimiter(mockStorage, 5, 10)
	clearStorage(mockStorage)
	r := setupRouter(rateLimiter)

	for i := 0; i < 5; i++ {
		fmt.Printf("Enviando requisição %d para o IP 192.168.1.1\n", i+1)
		req, _ := http.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		fmt.Printf("Requisição %d retornou status %d\n", i+1, rr.Code)
		assert.Equal(t, http.StatusOK, rr.Code, "Request should be allowed")
	}

	// Excedendo o limite
	fmt.Println("Enviando requisição que excede o limite para o IP 192.168.1.1")
	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	fmt.Printf("Requisição que excede o limite retornou status %d\n", rr.Code)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Request should be rate limited")
}

func TestRateLimitByToken(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	rateLimiter := ratelimiter.NewRateLimiter(mockStorage, 5, 10)
	clearStorage(mockStorage)
	r := setupRouter(rateLimiter)

	token := "abc123"

	for i := 0; i < 10; i++ {
		fmt.Printf("Enviando requisição %d com o token abc123\n", i+1)
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", token)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		fmt.Printf("Requisição %d retornou status %d\n", i+1, rr.Code)
		assert.Equal(t, http.StatusOK, rr.Code, "Request should be allowed")
	}

	// Excedendo o limite
	fmt.Println("Enviando requisição que excede o limite com o token abc123")
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	fmt.Printf("Requisição que excede o limite retornou status %d\n", rr.Code)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Request should be rate limited")
}

func TestRateLimitReset(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	rateLimiter := ratelimiter.NewRateLimiter(mockStorage, 5, 10)
	clearStorage(mockStorage)
	r := setupRouter(rateLimiter)

	clientIP := "192.168.1.2"

	for i := 0; i < 5; i++ {
		fmt.Printf("Enviando requisição %d para o IP 192.168.1.2\n", i+1)
		req, _ := http.NewRequest("GET", "/", nil)
		req.RemoteAddr = clientIP + ":1234"
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		fmt.Printf("Requisição %d retornou status %d\n", i+1, rr.Code)
		assert.Equal(t, http.StatusOK, rr.Code, "Request should be allowed")
	}

	// Excedendo o limite
	fmt.Println("Enviando requisição que excede o limite para o IP 192.168.1.2")
	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = clientIP + ":1234"
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	fmt.Printf("Requisição que excede o limite retornou status %d\n", rr.Code)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Request should be rate limited")

	// Aguardando o reset do limite
	fmt.Println("Aguardando 60 segundos para o reset do limite")
	time.Sleep(60 * time.Second)

	fmt.Println("Enviando requisição após o reset do limite para o IP 192.168.1.2")
	req, _ = http.NewRequest("GET", "/", nil)
	req.RemoteAddr = clientIP + ":1234"
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	fmt.Printf("Requisição após o reset retornou status %d\n", rr.Code)
	assert.Equal(t, http.StatusOK, rr.Code, "Request should be allowed after reset")
}

func TestRateLimitByTokenAndIP(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	rateLimiter := ratelimiter.NewRateLimiter(mockStorage, 5, 10)
	clearStorage(mockStorage)
	r := setupRouter(rateLimiter)

	token := "def456"
	clientIP := "192.168.1.3"

	for i := 0; i < 10; i++ {
		fmt.Printf("Enviando requisição %d com o token def456\n", i+1)
		req, _ := http.NewRequest("GET", "/", nil)
		req.RemoteAddr = clientIP + ":1234"
		req.Header.Set("API_KEY", token)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		fmt.Printf("Requisição %d retornou status %d\n", i+1, rr.Code)
		assert.Equal(t, http.StatusOK, rr.Code, "Request should be allowed")
	}

	// Excedendo o limite do token
	fmt.Println("Enviando requisição que excede o limite com o token def456")
	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = clientIP + ":1234"
	req.Header.Set("API_KEY", token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	fmt.Printf("Requisição que excede o limite retornou status %d\n", rr.Code)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code, "Request should be rate limited")

	// Verificando se IP ainda pode fazer requisições
	for i := 0; i < 5; i++ {
		fmt.Printf("Enviando requisição %d para o IP 192.168.1.3 sem o token\n", i+1)
		req, _ := http.NewRequest("GET", "/", nil)
		req.RemoteAddr = clientIP + ":1234"
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		fmt.Printf("Requisição %d retornou status %d\n", i+1, rr.Code)
		assert.Equal(t, http.StatusOK, rr.Code, "Request should be allowed")
	}
}
