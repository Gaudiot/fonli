package storytranslation

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"gaudiot.com/fonli/base"
	"gaudiot.com/fonli/core/analytics"
	"gaudiot.com/fonli/core/middlewares"
	"github.com/go-chi/chi/v5"
)

// MARK: - Payloads

type evaluateTranslationRequest struct {
	Story           string `json:"story"`
	UserTranslation string `json:"userTranslation"`
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

func StoryTranslationRouter(st *StoryTranslation) func(chi.Router) {
	return func(router chi.Router) {
		router.Get("/generate", handleGenerateStory(st))
		router.Post("/evaluate", handleEvaluateTranslation(st))
	}
}

// MARK: - Handlers

func handleGenerateStory(st *StoryTranslation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nativeLanguageCode := r.URL.Query().Get("nl")
		foreignLanguageCode := r.URL.Query().Get("fl")
		userID, _ := middlewares.UserIDFromContext(r.Context())

		if base.LanguageFromCountryCode(nativeLanguageCode) == "" || base.LanguageFromCountryCode(foreignLanguageCode) == "" {
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseStoryTranslationGenerate, analytics.ExerciseOutcomeValidationError,
				errors.New("invalid language code for 'nl' or 'fl'"))
			writeError(w, http.StatusBadRequest, "invalid language code for 'nl' or 'fl'")
			return
		}

		if userID == "" {
			writeError(w, http.StatusUnauthorized, "missing user id")
			return
		}

		story, err := st.GenerateStory(nativeLanguageCode, foreignLanguageCode, userID)
		if err != nil {
			slog.Error("failed to generate story", "error", err)
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseStoryTranslationGenerate, analytics.ExerciseOutcomeInternalError, err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		analytics.TrackExerciseInvocation(userID, analytics.ExerciseStoryTranslationGenerate, analytics.ExerciseOutcomeSuccess)
		writeJSON(w, http.StatusOK, story)
	}
}

const maxStoryLength = 5000

func handleEvaluateTranslation(st *StoryTranslation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nativeLanguageCode := r.URL.Query().Get("nl")
		foreignLanguageCode := r.URL.Query().Get("fl")
		userID, _ := middlewares.UserIDFromContext(r.Context())

		if base.LanguageFromCountryCode(nativeLanguageCode) == "" || base.LanguageFromCountryCode(foreignLanguageCode) == "" {
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseStoryTranslationEvaluate, analytics.ExerciseOutcomeValidationError,
				errors.New("invalid language code for 'nl' or 'fl'"))
			writeError(w, http.StatusBadRequest, "invalid language code for 'nl' or 'fl'")
			return
		}

		var req evaluateTranslationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Warn("invalid request body", "error", err)
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseStoryTranslationEvaluate, analytics.ExerciseOutcomeValidationError, err)
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if len([]rune(req.Story)) > maxStoryLength {
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseStoryTranslationEvaluate, analytics.ExerciseOutcomeValidationError,
				errors.New("story exceeds maximum length of 5000 characters"))
			writeError(w, http.StatusBadRequest, "story exceeds maximum length of 5000 characters")
			return
		}
		if len([]rune(req.UserTranslation)) > maxStoryLength {
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseStoryTranslationEvaluate, analytics.ExerciseOutcomeValidationError,
				errors.New("userTranslation exceeds maximum length of 5000 characters"))
			writeError(w, http.StatusBadRequest, "userTranslation exceeds maximum length of 5000 characters")
			return
		}

		if userID == "" {
			writeError(w, http.StatusUnauthorized, "missing user id")
			return
		}

		evaluation, err := st.EvaluateTranslation(req.Story, req.UserTranslation, nativeLanguageCode, foreignLanguageCode)
		if err != nil {
			slog.Error("failed to evaluate translation", "error", err)
			analytics.TrackExerciseInvocation(userID, analytics.ExerciseStoryTranslationEvaluate, analytics.ExerciseOutcomeInternalError, err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		analytics.TrackExerciseInvocation(userID, analytics.ExerciseStoryTranslationEvaluate, analytics.ExerciseOutcomeSuccess)
		writeJSON(w, http.StatusOK, evaluation)
	}
}
