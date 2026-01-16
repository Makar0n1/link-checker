package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/link-tracker/auth-service/internal/middleware"
	"github.com/link-tracker/auth-service/internal/model"
	"github.com/link-tracker/auth-service/internal/repository"
	"github.com/link-tracker/auth-service/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	if err := validateRegisterRequest(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			respondError(w, http.StatusConflict, "user with this email already exists", "USER_EXISTS")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to register user", "INTERNAL_ERROR")
		return
	}

	respondJSON(w, http.StatusCreated, model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// Login handles user authentication
// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "email and password are required", "VALIDATION_ERROR")
		return
	}

	authResponse, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			respondError(w, http.StatusUnauthorized, "invalid email or password", "INVALID_CREDENTIALS")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to login", "INTERNAL_ERROR")
		return
	}

	respondJSON(w, http.StatusOK, authResponse)
}

// Refresh handles token refresh
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	if req.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "refresh_token is required", "VALIDATION_ERROR")
		return
	}

	authResponse, err := h.authService.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			respondError(w, http.StatusUnauthorized, "invalid refresh token", "INVALID_TOKEN")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to refresh tokens", "INTERNAL_ERROR")
		return
	}

	respondJSON(w, http.StatusOK, authResponse)
}

// Me returns current user information
// GET /api/v1/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			respondError(w, http.StatusNotFound, "user not found", "USER_NOT_FOUND")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get user", "INTERNAL_ERROR")
		return
	}

	respondJSON(w, http.StatusOK, model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// Logout handles user logout
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	if req.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "refresh_token is required", "VALIDATION_ERROR")
		return
	}

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to logout", "INTERNAL_ERROR")
		return
	}

	respondJSON(w, http.StatusOK, model.MessageResponse{
		Message: "successfully logged out",
	})
}

func validateRegisterRequest(req *model.RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.ErrorResponse{
		Error: message,
		Code:  code,
	})
}
