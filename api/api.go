package api

import (
	"api/omdb"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Response struct {
	Error string      `json:"error,omitempty"`
	Data  omdb.Result `json:"data,omitempty"`
}

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json", "error", err)
		sendJSON(w, Response{Error: "something went wrong"}, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response to client", "error", err)
		return
	}
}

func NewHandler(apiKey string) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/", handleSearchMovie(apiKey))

	return r
}

func handleSearchMovie(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("s")
		res, err := omdb.Search(apiKey, title)
		if err != nil {
			sendJSON(w, Response{Error: "something went wrong with omdb"}, http.StatusBadGateway)
			return
		}

		sendJSON(w, Response{Data: res}, http.StatusOK)
	}
}
