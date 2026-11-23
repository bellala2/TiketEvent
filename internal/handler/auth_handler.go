package handler

import (
	"context"
	"net/http"

	"github.com/bellala2/TiketEvent/internal/middleware"
	"github.com/bellala2/TiketEvent/internal/model"
	"github.com/bellala2/TiketEvent/internal/response"
	"github.com/bellala2/TiketEvent/internal/service"
)

type authService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.RegisteredUser, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	GetMe(ctx context.Context, userID string) (*model.UserResponse, error)
}

type AuthHandler struct {
	svc authService
}

func NewAuthHandler(svc authService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register menjawab POST /auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := decodeJSON(w, r, &req); err != nil {
		response.WriteError(w, err)
		return
	}

	user, err := h.svc.Register(r.Context(), req)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteSuccess(w, http.StatusCreated, "Registrasi berhasil", user, nil)
}

// Login menjawab POST /auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := decodeJSON(w, r, &req); err != nil {
		response.WriteError(w, err)
		return
	}

	result, err := h.svc.Login(r.Context(), req)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Login berhasil", result, nil)
}

// Me menjawab GET /auth/me. Harus lewat middleware Authenticate.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, service.NewError(service.CodeUnauthorized, "Token tidak valid atau sudah kedaluwarsa"))
		return
	}

	user, err := h.svc.GetMe(r.Context(), authUser.ID)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteSuccess(w, http.StatusOK, "Data user berhasil diambil", user, nil)
}
