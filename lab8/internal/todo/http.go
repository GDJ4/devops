package todo

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	repo *Repository
}

func NewServer(repo *Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/todos", s.handleTodos)
	mux.HandleFunc("/todos/", s.handleTodoByID)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		mux.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		todos, err := s.repo.List(ctx)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, todos)
	case http.MethodPost:
		var req CreateTodoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, errors.New("invalid JSON payload"))
			return
		}
		req.Title = strings.TrimSpace(req.Title)
		if req.Title == "" {
			writeError(w, http.StatusBadRequest, errors.New("title is required"))
			return
		}
		todo, err := s.repo.Create(ctx, req)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusCreated, todo)
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) handleTodoByID(w http.ResponseWriter, r *http.Request) {
	idHex := strings.TrimPrefix(r.URL.Path, "/todos/")
	if idHex == "" {
		writeError(w, http.StatusNotFound, errors.New("missing todo id"))
		return
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid id format"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		todo, err := s.repo.Get(ctx, id)
		if err != nil {
			s.handleRepoError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, todo)
	case http.MethodPut:
		var req UpdateTodoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, errors.New("invalid JSON payload"))
			return
		}
		if req.Title == nil && req.Completed == nil && req.Notes == nil {
			writeError(w, http.StatusBadRequest, errors.New("no fields to update"))
			return
		}
		if req.Title != nil {
			trimmed := strings.TrimSpace(*req.Title)
			req.Title = &trimmed
		}
		todo, err := s.repo.Update(ctx, id, req)
		if err != nil {
			s.handleRepoError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, todo)
	case http.MethodDelete:
		if err := s.repo.Delete(ctx, id); err != nil {
			s.handleRepoError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) handleRepoError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func methodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
}
