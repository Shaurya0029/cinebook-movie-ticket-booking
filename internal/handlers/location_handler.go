package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"movieticketbooking/internal/geo"
	"movieticketbooking/internal/httpx"
)

func (d *Deps) ListCities(w http.ResponseWriter, r *http.Request) {
	cities, err := d.Cities.List(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list cities")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, cities)
}

func (d *Deps) NearestCity(w http.ResponseWriter, r *http.Request) {
	lat, lng, ok := parseLatLng(r)
	if !ok {
		httpx.WriteError(w, http.StatusBadRequest, "lat and lng query params are required")
		return
	}

	cities, err := d.Cities.List(r.Context())
	if err != nil || len(cities) == 0 {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list cities")
		return
	}

	best := cities[0]
	bestDist := geo.DistanceKm(lat, lng, best.Latitude, best.Longitude)
	for _, c := range cities[1:] {
		dist := geo.DistanceKm(lat, lng, c.Latitude, c.Longitude)
		if dist < bestDist {
			best, bestDist = c, dist
		}
	}
	httpx.WriteJSON(w, http.StatusOK, best)
}

func (d *Deps) ListTheatersByCity(w http.ResponseWriter, r *http.Request) {
	cityID, err := strconv.ParseInt(chi.URLParam(r, "cityID"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid city id")
		return
	}
	theaters, err := d.Theaters.ListByCity(r.Context(), cityID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list theaters")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, theaters)
}

func (d *Deps) GetTheater(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "theaterID"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid theater id")
		return
	}
	theater, err := d.Theaters.GetByID(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, "theater not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, theater)
}

func (d *Deps) NearestTheater(w http.ResponseWriter, r *http.Request) {
	lat, lng, ok := parseLatLng(r)
	if !ok {
		httpx.WriteError(w, http.StatusBadRequest, "lat and lng query params are required")
		return
	}

	theaters, err := d.Theaters.ListAllWithCity(r.Context())
	if err != nil || len(theaters) == 0 {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list theaters")
		return
	}

	best := theaters[0]
	bestDist := geo.DistanceKm(lat, lng, best.Latitude, best.Longitude)
	for _, t := range theaters[1:] {
		dist := geo.DistanceKm(lat, lng, t.Latitude, t.Longitude)
		if dist < bestDist {
			best, bestDist = t, dist
		}
	}
	httpx.WriteJSON(w, http.StatusOK, best)
}

func parseLatLng(r *http.Request) (float64, float64, bool) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")
	lat, err1 := strconv.ParseFloat(latStr, 64)
	lng, err2 := strconv.ParseFloat(lngStr, 64)
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return lat, lng, true
}
