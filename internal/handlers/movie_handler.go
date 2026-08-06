package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"movieticketbooking/internal/httpx"
)

func (d *Deps) ListMovies(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status != "" && status != "now_showing" && status != "coming_soon" {
		httpx.WriteError(w, http.StatusBadRequest, "status must be now_showing or coming_soon")
		return
	}
	movies, err := d.Movies.List(r.Context(), status)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list movies")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, movies)
}

func (d *Deps) GetMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "movieID"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid movie id")
		return
	}
	movie, err := d.Movies.GetByID(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, "movie not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, movie)
}
