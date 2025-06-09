package repo

import (
	"backend/internal/model"
	"sync"
)

type Repository interface {
	GetAll() []model.Tournament
	Save(model.Tournament) model.Tournament
	Delete(id uint)
	Archive(id uint)
}

type InMemoryRepo struct {
	mu          sync.RWMutex
	tournaments map[uint]model.Tournament
	nextID      uint
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		tournaments: make(map[uint]model.Tournament),
		nextID:      1,
	}
}

func (r *InMemoryRepo) GetAll() []model.Tournament {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]model.Tournament, 0, len(r.tournaments))
	for _, t := range r.tournaments {
		list = append(list, t)
	}
	return list
}

func (r *InMemoryRepo) Save(t model.Tournament) model.Tournament {
	r.mu.Lock()
	defer r.mu.Unlock()
	t.ID = r.nextID
	r.nextID++
	r.tournaments[t.ID] = t
	return t
}

func (r *InMemoryRepo) Delete(id uint) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tournaments, id)
}

func (r *InMemoryRepo) Archive(id uint) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t, ok := r.tournaments[id]; ok {
		t.Status = model.Archive
		r.tournaments[id] = t
	}
}
