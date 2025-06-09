package api

import (
	"backend/internal/model"
	"backend/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	svc    *service.TournamentService
	logger *zap.Logger
}

func NewHandler(svc *service.TournamentService, logger *zap.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/tournaments", h.GetAll).Methods("GET")
	r.HandleFunc("/api/tournaments", h.Create).Methods("POST")
	r.HandleFunc("/api/tournaments/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/api/tournaments/{id}", h.ChangeState).Methods("PATCH")
}

func (h *Handler) GetAll(w http.ResponseWriter, _ *http.Request) {
	list := h.svc.GetAll()
	h.logger.Info("get", zap.Any("saved", list))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var t model.Tournament
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.logger.Error("invalid create payload", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if t.Date.IsZero() {
		t.Date = time.Now().UTC()
	}
	saved := h.svc.Create(t)
	h.logger.Info("saved", zap.Any("saved", saved))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saved)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	h.svc.Delete(uint(id))
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ChangeState(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	state := model.TournamentState(r.URL.Query().Get("state"))
	if state != model.Active && state != model.Archive && state != model.InProgress {
		h.logger.Error("invalid change payload")
		http.Error(w, "invalid change payload", http.StatusBadRequest)
		return
	}
	h.svc.ChangeState(uint(id), state)
	w.WriteHeader(http.StatusNoContent)
}
