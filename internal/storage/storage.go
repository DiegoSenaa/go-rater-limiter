package storage

type Storage interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration int) error
	Incr(key string) error
	Clear() error
}
