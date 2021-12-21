package db

type InMemoryStore interface {
	Get(string) (string, error)
	Set(string, string) error
}
