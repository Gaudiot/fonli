package wordconjugationexercise

import (
	"encoding/json"
	"net/http"

	"gaudiot.com/fonli/base"
	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	"github.com/go-chi/chi/v5"
)

func WordConjugationRouter(router chi.Router) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		aiService := aiservice.NewOpenAIAIService()
		wordConjugation := NewWordConjugation(aiService)

		query := r.URL.Query()
		rawTense := query.Get("tense")
		tense := base.GetTense(rawTense)

		exercises, err := wordConjugation.GenerateExercise(tense)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercises)
	})
}
