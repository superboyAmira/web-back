package model

import "time"

type Tournament struct {
	ID        uint            `json:"id"`
	Title     string          `json:"title"`
	Date      time.Time       `json:"date"`
	Location  string          `json:"location"`
	Status    TournamentState `json:"status"`
	Organizer string          `json:"organizer"`
}

type TournamentState string

const (
	Active     TournamentState = "active"
	Archive    TournamentState = "archive"
	InProgress TournamentState = "in_progress"
)

func (t TournamentState) ToString() string {
	return string(t)
}

/*
{
  "id": 1,
  "title": "Московский блиц-турнир",
  "date": "2025-06-25T18:00:00",
  "location": "Москва, Центральный шахматный клуб",
  "status": "active", // или "archive"
  "organizer": "Федерация шахмат РФ"
}
*/
