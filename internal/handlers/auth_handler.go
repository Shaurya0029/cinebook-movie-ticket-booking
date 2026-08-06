package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"movieticketbooking/internal/auth"
	"movieticketbooking/internal/httpx"
	"movieticketbooking/internal/repository"
)

type registerRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	User  userResponse `json:"user"`
	Token string       `json:"token"`
}

type userResponse struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (d *Deps) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if len(req.FirstName) < 2 || len(req.LastName) < 2 {
		httpx.WriteError(w, http.StatusBadRequest, "first and last name must be at least 2 characters")
		return
	}
	if !strings.Contains(req.Email, "@") {
		httpx.WriteError(w, http.StatusBadRequest, "invalid email address")
		return
	}
	if len(req.Password) < 8 {
		httpx.WriteError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := d.Users.Create(r.Context(), req.FirstName, req.LastName, req.Email, hash)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			httpx.WriteError(w, http.StatusConflict, "email already registered")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	token, err := d.TokenIssuer.Issue(user.ID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to issue token")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, authResponse{
		User:  userResponse{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email},
		Token: token,
	})
}

func (d *Deps) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	user, err := d.Users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := d.TokenIssuer.Issue(user.ID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to issue token")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, authResponse{
		User:  userResponse{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email},
		Token: token,
	})
}

func (d *Deps) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	user, err := d.Users.GetByID(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, userResponse{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email})
}
