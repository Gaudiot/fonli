package wordtranslationexercise

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

func WordTranslationRouter(wt *WordTranslation) func(chi.Router) {
	return func(router chi.Router) {
		router.Get("/native-to-foreign", handleNativeToForeignExercise(wt))
		router.Get("/foreign-to-native", handleForeignToNativeExercise(wt))
	}
}

// MARK: - Handlers

func handleNativeToForeignExercise(wt *WordTranslation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nl := r.URL.Query().Get("nl")
		fl := r.URL.Query().Get("fl")

		if base.LanguageFromCountryCode(nl) == "" || base.LanguageFromCountryCode(fl) == "" {
			writeError(w, http.StatusBadRequest, "invalid language code for 'nl' or 'fl'")
			return
		}

		userID, ok := middlewares.UserIDFromContext(r.Context())
		if !ok || userID == "" {
			writeError(w, http.StatusUnauthorized, "missing user id")
			return
		}

		exercises, err := wt.NativeToForeignExercise(10, nl, fl, userID)
		if err != nil {
			slog.Error("failed to generate native-to-foreign exercise", "error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusOK, exercises)
	}
}

func handleForeignToNativeExercise(wt *WordTranslation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nl := r.URL.Query().Get("nl")
		fl := r.URL.Query().Get("fl")

		if base.LanguageFromCountryCode(nl) == "" || base.LanguageFromCountryCode(fl) == "" {
			writeError(w, http.StatusBadRequest, "invalid language code for 'nl' or 'fl'")
			return
		}

		userID, ok := middlewares.UserIDFromContext(r.Context())
		if !ok || userID == "" {
			writeError(w, http.StatusUnauthorized, "missing user id")
			return
		}

		exercises, err := wt.ForeignToNativeExercise(10, fl, nl, userID)
		if err != nil {
			slog.Error("failed to generate foreign-to-native exercise", "error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusOK, exercises)
	}
}
