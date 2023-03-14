package core

type Storage interface {
	Put(*Block) error
}

type MemoryStore struct{}

func NewMemoryStorage() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Put(*Block) error {
	return nil
}
