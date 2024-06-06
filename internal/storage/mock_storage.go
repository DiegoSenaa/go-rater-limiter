package storage

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type MockStorage struct {
	data map[string]string
	mu   sync.Mutex
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		data: make(map[string]string),
	}
}

func (m *MockStorage) Get(key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return "", redis.Nil
}

func (m *MockStorage) Set(key string, value string, expiration int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	go func() {
		time.Sleep(time.Duration(expiration) * time.Second)
		m.mu.Lock()
		delete(m.data, key)
		m.mu.Unlock()
	}()
	return nil
}

func (m *MockStorage) Incr(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if val, ok := m.data[key]; ok {
		num, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		m.data[key] = strconv.Itoa(num + 1)
		return nil
	}
	return errors.New("key not found")
}

func (m *MockStorage) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]string)
	return nil
}
