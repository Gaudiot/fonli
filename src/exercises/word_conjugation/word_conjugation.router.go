package wordconjugationexercise

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"gaudiot.com/fonli/base"
	"gaudiot.com/fonli/core/middlewares"
	"github.com/go-chi/chi/v5"
)

// MARK: - Helpers

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

// MARK: - Router

func WordConjugationRouter(wc *WordConjugation) func(chi.Router) {
	return func(router chi.Router) {
		router.Get("/", handleGenerateExercise(wc))
	}
}

// MARK: - Handlers

func handleGenerateExercise(wc *WordConjugation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foreignLanguageCode := r.URL.Query().Get("fl")
		rawTense := r.URL.Query().Get("tense")

		if base.LanguageFromCountryCode(foreignLanguageCode) == "" {
			writeError(w, http.StatusBadRequest, "invalid language code for 'fl'")
			return
		}

		if rawTense == "" {
			writeError(w, http.StatusBadRequest, "'tense' query parameter is required")
			return
		}

		tense := base.GetTense(rawTense)

		userID, ok := middlewares.UserIDFromContext(r.Context())
		if !ok || userID == "" {
			writeError(w, http.StatusUnauthorized, "missing user id")
			return
		}

		exercises, err := wc.GenerateExercise(tense, foreignLanguageCode, userID)
		if err != nil {
			slog.Error("failed to generate conjugation exercise", "error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusOK, exercises)
	}
}
