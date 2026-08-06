package handlers

import (
	"net/http"

	"movieticketbooking/internal/httpx"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
