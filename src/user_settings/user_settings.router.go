package user_settings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gaudiot.com/fonli/core/middlewares"
	"gaudiot.com/fonli/core/security/tokens"
	"github.com/go-chi/chi/v5"
)

// MARK: - Payloads

type updateLifestyleRequest struct {
	Text string `json:"text"`
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

// MARK: - Router

func UserSettingsRouter(us *UserSettingsService, ts *tokens.TokenService) func(chi.Router) {
	return func(router chi.Router) {
		router.Use(middlewares.AuthMiddleware(ts))
		router.Get("/lifestyle", handleGetLifestyle(us))
		router.Post("/lifestyle", handleUpdateLifestyle(us))
	}
}

// MARK: - Handlers

func handleGetLifestyle(us *UserSettingsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userExists := middlewares.UserIDFromContext(ctx)
		if !userExists {
			fmt.Println("user not authenticated")
			writeError(w, http.StatusUnauthorized, "user not authenticated")
			return
		}

		lifestyle, err := us.GetUserLifestyle(userID)
		if err != nil {
			fmt.Println("internal server error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(w, http.StatusOK, lifestyle)
	}
}

func handleUpdateLifestyle(us *UserSettingsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, userExists := middlewares.UserIDFromContext(ctx)
		if !userExists {
			fmt.Println("user not authenticated")
			writeError(w, http.StatusUnauthorized, "user not authenticated")
			return
		}

		var req updateLifestyleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Println("invalid request body", err)
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		err := us.UpdateUserLifestyle(userID, req.Text)
		if err != nil {
			fmt.Println("internal server error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		fmt.Println("lifestyle updated successfully")
		writeJSON(w, http.StatusOK, "lifestyle updated successfully")
	}
}
