package wordtranslationexercise

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"gaudiot.com/fonli/base"
	"gaudiot.com/fonli/core/analytics"
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
		baseLanguageCode := r.URL.Query().Get("nl")
		targetLanguageCode := r.URL.Query().Get("fl")
		userID, _ := middlewares.UserIDFromContext(r.Context())

		baseLanguage, err := base.LanguageFromCountryCode(baseLanguageCode)
		if err != nil {
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseWordTranslationNativeToForeign, analytics.ExerciseOutcomeValidationError, err)
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		targetLanguage, err := base.LanguageFromCountryCode(targetLanguageCode)
		if err != nil {
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseWordTranslationNativeToForeign, analytics.ExerciseOutcomeValidationError, err)
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		if userID == "" {
			writeError(w, http.StatusUnauthorized, "missing user id")
			return
		}

		exercises, err := wt.NativeToForeignExercise(10, baseLanguage, targetLanguage, userID)
		if err != nil {
			slog.Error("failed to generate native-to-foreign exercise", "error", err)
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseWordTranslationNativeToForeign, analytics.ExerciseOutcomeInternalError, err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		analytics.TrackExerciseInvocation(userID, analytics.ExerciseWordTranslationNativeToForeign, analytics.ExerciseOutcomeSuccess)
		writeJSON(w, http.StatusOK, exercises)
	}
}

func handleForeignToNativeExercise(wt *WordTranslation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseLanguageCode := r.URL.Query().Get("nl")
		targetLanguageCode := r.URL.Query().Get("fl")
		userID, _ := middlewares.UserIDFromContext(r.Context())

		baseLanguage, err := base.LanguageFromCountryCode(baseLanguageCode)
		if err != nil {
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseWordTranslationForeignToNative, analytics.ExerciseOutcomeValidationError, err)
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		targetLanguage, err := base.LanguageFromCountryCode(targetLanguageCode)
		if err != nil {
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseWordTranslationForeignToNative, analytics.ExerciseOutcomeValidationError, err)
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		if userID == "" {
			writeError(w, http.StatusUnauthorized, "missing user id")
			return
		}

		exercises, err := wt.ForeignToNativeExercise(10, baseLanguage, targetLanguage, userID)
		if err != nil {
			slog.Error("failed to generate foreign-to-native exercise", "error", err)
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseWordTranslationForeignToNative, analytics.ExerciseOutcomeInternalError, err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		analytics.TrackExerciseInvocation(userID, analytics.ExerciseWordTranslationForeignToNative, analytics.ExerciseOutcomeSuccess)
		writeJSON(w, http.StatusOK, exercises)
	}
}
