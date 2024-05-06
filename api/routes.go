package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"words/repository"
)

// Interface vers PG store
type handler struct {
	repo repository.Store
}

func NewHandler(repo repository.Store) *handler {
	return &handler{repo}
}

type errorResponse struct {
	Prefix string `json:"prefix"`
	Error  string `json:"error"`
}

func writeError(cause errorResponse, w http.ResponseWriter, status int) {
	data, err := json.Marshal(cause)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(status)
	w.Write(data)
}

// POST /service/{word}
func (h *handler) InsertWord(w http.ResponseWriter, r *http.Request) {
	word := strings.ToLower(r.PathValue("word"))
	if !IsValid(word) {
		writeError(errorResponse{
			Prefix: word,
			Error:  "bad string format",
		}, w, http.StatusBadRequest)
		return
	}
	if err := h.repo.Insert(r.Context(), word); err != nil {
		writeError(errorResponse{
			Prefix: word,
			Error:  err.Error(),
		}, w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GET /service/{prefix}
func (h *handler) FindPrefix(w http.ResponseWriter, r *http.Request) {
	prefix := strings.ToLower(r.PathValue("prefix"))
	if !IsValid(prefix) {
		writeError(errorResponse{
			Prefix: prefix,
			Error:  "bad string format",
		}, w, http.StatusBadRequest)
		return
	}
	result, err := h.repo.GetByPrefix(r.Context(), prefix)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(errorResponse{
			Prefix: prefix,
			Error:  err.Error(),
		}, w, status)
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	words, err := h.repo.List(r.Context())
	if err != nil {
		writeError(errorResponse{
			Prefix: "",
			Error:  err.Error(),
		}, w, http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(words)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

var re = regexp.MustCompile(`^[a-zA-Z]+$`)

func IsValid(w string) bool {
	if len(w) > 32 {
		return false
	}
	return re.MatchString(w)
}
