package wordconjugationexercise

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gaudiot.com/fonli/base"
	"gaudiot.com/fonli/core/middlewares"
	"gaudiot.com/fonli/core/security/tokens"
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

func WordConjugationRouter(wc *WordConjugation, ts *tokens.TokenService) func(chi.Router) {
	return func(router chi.Router) {
		router.Use(middlewares.AuthMiddleware(ts))
		router.Get("/", handleGenerateExercise(wc))
	}
}

// MARK: - Handlers

func handleGenerateExercise(wc *WordConjugation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foreignLanguageCode := r.URL.Query().Get("fl")
		rawTense := r.URL.Query().Get("tense")
		tense := base.GetTense(rawTense)

		exercises, err := wc.GenerateExercise(tense, foreignLanguageCode)
		if err != nil {
			fmt.Println("internal server error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusOK, exercises)
	}
}
