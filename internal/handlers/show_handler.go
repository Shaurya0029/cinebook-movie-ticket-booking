package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"movieticketbooking/internal/httpx"
	"movieticketbooking/internal/repository"
)

func (d *Deps) ListShowsForMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.ParseInt(chi.URLParam(r, "movieID"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid movie id")
		return
	}

	filter := repository.ShowFilter{MovieID: movieID}

	if v := r.URL.Query().Get("theater_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid theater_id")
			return
		}
		filter.TheaterID = id
	}
	if v := r.URL.Query().Get("city_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid city_id")
			return
		}
		filter.CityID = id
	}
	if v := r.URL.Query().Get("date"); v != "" {
		date, err := time.Parse("2006-01-02", v)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid date, expected YYYY-MM-DD")
			return
		}
		filter.Date = &date
	}

	shows, err := d.Shows.ListForMovie(r.Context(), filter)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list shows")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, shows)
}

func (d *Deps) GetShow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "showID"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid show id")
		return
	}
	show, err := d.Shows.GetByID(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, "show not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, show)
}

func (d *Deps) ListSeatsForShow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "showID"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid show id")
		return
	}
	seats, err := d.Seats.ListForShow(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list seats")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, seats)
}
