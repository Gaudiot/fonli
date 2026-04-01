package storytranslation

import (
	"encoding/json"
	"fmt"
	"net/http"

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

		story, err := st.GenerateStory(nativeLanguageCode, foreignLanguageCode)
		if err != nil {
			fmt.Println("internal server error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusOK, story)
	}
}

func handleEvaluateTranslation(st *StoryTranslation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nativeLanguageCode := r.URL.Query().Get("nl")
		foreignLanguageCode := r.URL.Query().Get("fl")

		var req evaluateTranslationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Println("invalid request body", err)
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		evaluation, err := st.EvaluateTranslation(req.Story, req.UserTranslation, nativeLanguageCode, foreignLanguageCode)
		if err != nil {
			fmt.Println("internal server error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusOK, evaluation)
	}
}
