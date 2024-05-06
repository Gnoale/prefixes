package api

import (
	"encoding/json"
	"net/http"
	"regexp"
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
	word := r.PathValue("word")
	if !IsValid(word) {
		writeError(errorResponse{
			Prefix: word,
			Error:  "bad string format",
		}, w, http.StatusInternalServerError)
		return
	}
	if err := h.repo.Insert(word); err != nil {
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
	prefix := r.PathValue("prefix")
	result, err := h.repo.GetByPrefix(prefix)
	if err != nil {
		writeError(errorResponse{
			Prefix: prefix,
			Error:  err.Error(),
		}, w, http.StatusInternalServerError)
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

var re = regexp.MustCompile(`[a-zA-Z]+`)

func IsValid(w string) bool {
	return re.MatchString(w)
}
