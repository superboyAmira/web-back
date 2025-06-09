package service

import "backend/internal/model"

type Repository interface {
	GetAll() []model.Tournament
	Save(model.Tournament) model.Tournament
	Delete(id uint)
	Archive(id uint)
}

type TournamentService struct {
	repo Repository
}

func NewTournamentService(r Repository) *TournamentService {
	return &TournamentService{repo: r}
}

func (s *TournamentService) GetAll() []model.Tournament {
	return s.repo.GetAll()
}

func (s *TournamentService) Create(t model.Tournament) model.Tournament {
	return s.repo.Save(t)
}

func (s *TournamentService) Delete(id uint) {
	s.repo.Delete(id)
}

func (s *TournamentService) Archive(id uint) {
	s.repo.Archive(id)
}
