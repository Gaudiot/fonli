package authentication

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// MARK: - Payloads

type signUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type authTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// MARK: - Helpers

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func mapSignUpError(err error) (int, string) {
	switch {
	case errors.Is(err, ErrInvalidEmail),
		errors.Is(err, ErrPasswordTooShort),
		errors.Is(err, ErrPasswordTooLong),
		errors.Is(err, ErrPasswordMissingLetter),
		errors.Is(err, ErrPasswordMissingDigit),
		errors.Is(err, ErrPasswordNoWhitespace),
		errors.Is(err, ErrUsernameTooShort),
		errors.Is(err, ErrUsernameTooLong),
		errors.Is(err, ErrUsernameInvalidFormat):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, ErrEmailAlreadyRegistered),
		errors.Is(err, ErrUsernameAlreadyRegistered):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

// MARK: - Router

func AuthenticationRouter(as *AuthService) func(chi.Router) {
	return func(router chi.Router) {
		router.Post("/signup", handleSignUp(as))
		router.Post("/login", handleLogin(as))
		router.Post("/refresh", handleRefresh(as))
		router.Post("/logout", handleLogout(as))
	}
}

// MARK: - Handlers

func handleSignUp(as *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		tokens, err := as.SignUp(req.Username, req.Email, req.Password)
		if err != nil {
			status, message := mapSignUpError(err)
			if status == http.StatusInternalServerError {
				slog.Error("signup failed", "error", err)
			}
			writeError(w, status, message)
			return
		}

		writeJSON(w, http.StatusCreated, authTokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}

func handleLogin(as *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		tokens, err := as.Login(req.EmailOrUsername, req.Password)
		if err != nil {
			if errors.Is(err, ErrInvalidCredentials) {
				writeError(w, http.StatusUnauthorized, err.Error())
				return
			}
			slog.Error("login failed", "error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusOK, authTokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}

func handleRefresh(as *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req refreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		tokens, err := as.Refresh(req.RefreshToken)
		if err != nil {
			if errors.Is(err, ErrInvalidRefreshToken) || errors.Is(err, ErrExpiredRefreshToken) {
				writeError(w, http.StatusUnauthorized, err.Error())
				return
			}
			slog.Error("refresh failed", "error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusOK, authTokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}

func handleLogout(as *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := as.tokenService.ParseAccessToken(tokenString)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid access token")
			return
		}

		if err := as.Logout(claims.UserID); err != nil {
			slog.Error("logout failed", "error", err, "userID", claims.UserID)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
