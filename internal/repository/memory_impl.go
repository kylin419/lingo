package repository

import(
	"sync"
)

type memoryRepo struct{
	data map[string][]string
	mu sync.RWMutex
}

func NewMemoryRepository() GroupRepository {
	return &memoryRepo{
		data: make(map[string][]string),
	}
}

func (m *memoryRepo)SaveLanguages(sourceID string,langs []string)error{
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[sourceID]=langs
	return nil
}

func (m *memoryRepo)GetLanguages(sourceID string)([]string,error){
	m.mu.RLock()
	defer m.mu.RUnlock()

	langs,exists:=m.data[sourceID]
	if !exists{
		return []string{},nil
	}
	return langs,nil
}