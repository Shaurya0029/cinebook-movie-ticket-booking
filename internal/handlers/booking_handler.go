package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"movieticketbooking/internal/auth"
	"movieticketbooking/internal/httpx"
	"movieticketbooking/internal/models"
	"movieticketbooking/internal/repository"
)

type createBookingRequest struct {
	ShowID  int64   `json:"show_id"`
	SeatIDs []int64 `json:"seat_ids"`
}

func (d *Deps) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ShowID == 0 || len(req.SeatIDs) == 0 {
		httpx.WriteError(w, http.StatusBadRequest, "show_id and at least one seat_id are required")
		return
	}

	show, err := d.Shows.GetByID(r.Context(), req.ShowID)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, "show not found")
		return
	}
	if !show.StartsAt.After(time.Now()) {
		httpx.WriteError(w, http.StatusConflict, "this show has already started and can no longer be booked")
		return
	}

	booking, err := d.Bookings.CreateBooking(r.Context(), userID, show.ID, req.SeatIDs, show.PriceCents, d.GSTRatePercent)
	if err != nil {
		if errors.Is(err, repository.ErrSeatUnavailable) {
			httpx.WriteError(w, http.StatusConflict, "one or more selected seats are no longer available")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create booking")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, booking)
}

func (d *Deps) ListMyBookings(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	bookings, err := d.Bookings.ListForUser(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list bookings")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, bookings)
}

func (d *Deps) GetBooking(w http.ResponseWriter, r *http.Request) {
	booking, ok := d.loadOwnedBooking(w, r)
	if !ok {
		return
	}
	httpx.WriteJSON(w, http.StatusOK, booking)
}

func (d *Deps) GetBookingReceipt(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	booking, ok := d.loadOwnedBooking(w, r)
	if !ok {
		return
	}

	show, err := d.Shows.GetByID(r.Context(), booking.ShowID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to load show")
		return
	}
	user, err := d.Users.GetByID(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to load user")
		return
	}
	seatLabels, err := d.Bookings.SeatLabelsForBooking(r.Context(), booking.ID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to load seats")
		return
	}
	payment, err := d.Bookings.GetPaymentForBooking(r.Context(), booking.ID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to load payment")
		return
	}

	receipt := models.Receipt{
		Booking:       *booking,
		Payment:       payment,
		MovieTitle:    show.MovieTitle,
		TheaterName:   show.TheaterName,
		CityName:      show.CityName,
		ShowStartsAt:  show.StartsAt,
		PriceCents:    show.PriceCents,
		UserFirstName: user.FirstName,
		UserLastName:  user.LastName,
		SeatLabels:    seatLabels,
	}
	httpx.WriteJSON(w, http.StatusOK, receipt)
}

const testCardFailNumber = "4000000000000002"

type payRequest struct {
	CardNumber string `json:"card_number"`
	CardName   string `json:"card_name"`
	Expiry     string `json:"expiry"`
	CVV        string `json:"cvv"`
}

func (d *Deps) PayForBooking(w http.ResponseWriter, r *http.Request) {
	booking, ok := d.loadOwnedBooking(w, r)
	if !ok {
		return
	}
	if booking.Status != "pending_payment" {
		httpx.WriteError(w, http.StatusConflict, "booking is not awaiting payment")
		return
	}

	var req payRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cardDigits := strings.ReplaceAll(req.CardNumber, " ", "")
	if len(cardDigits) < 12 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid card number")
		return
	}
	last4 := cardDigits[len(cardDigits)-4:]

	// Simulated payment processor: a specific magic test card always fails,
	// everything else succeeds. No real payment gateway is involved.
	if cardDigits == testCardFailNumber {
		payment, err := d.Bookings.FailBooking(r.Context(), booking.ID, booking.TotalAmountCents, last4, "card declined (simulated)")
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "failed to record payment failure")
			return
		}
		httpx.WriteJSON(w, http.StatusPaymentRequired, map[string]any{"payment": payment, "booking_status": "cancelled"})
		return
	}

	payment, err := d.Bookings.ConfirmBooking(r.Context(), booking.ID, booking.TotalAmountCents, last4)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to confirm payment")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"payment": payment, "booking_status": "confirmed"})
}

// loadOwnedBooking loads the booking named in the URL and verifies it belongs
// to the authenticated user, writing an error response and returning ok=false otherwise.
func (d *Deps) loadOwnedBooking(w http.ResponseWriter, r *http.Request) (*models.Booking, bool) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "not authenticated")
		return nil, false
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "bookingID"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return nil, false
	}

	booking, err := d.Bookings.GetByID(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, "booking not found")
		return nil, false
	}
	if booking.UserID != userID {
		httpx.WriteError(w, http.StatusForbidden, "not your booking")
		return nil, false
	}
	return booking, true
}
