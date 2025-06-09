package repo

import (
	"backend/internal/model"
	"sync"
	"time"
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
	r := &InMemoryRepo{
		tournaments: make(map[uint]model.Tournament),
		nextID:      1,
	}

	r.Save(model.Tournament{
		Title:     "Московский блиц-турнир",
		Date:      time.Date(2025, 6, 25, 18, 0, 0, 0, time.UTC),
		Location:  "Москва, Центральный шахматный клуб",
		Organizer: "Федерация шахмат РФ",
		Status:    model.Active,
	})
	r.Save(model.Tournament{
		Title:     "Кубок чемпионов",
		Date:      time.Date(2025, 8, 1, 15, 30, 0, 0, time.UTC),
		Location:  "Санкт-Петербург, КСШ",
		Organizer: "Союз шахмат России",
		Status:    model.InProgress,
	})
	r.Save(model.Tournament{
		Title:     "Летний шахматный опен",
		Date:      time.Date(2025, 7, 10, 12, 0, 0, 0, time.UTC),
		Location:  "Сочи, Олимпийская деревня",
		Organizer: "Министерство спорта",
		Status:    model.Archive,
	})

	return r
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
