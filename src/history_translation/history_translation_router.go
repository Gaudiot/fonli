package historytranslation

import (
	"encoding/json"
	"net/http"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	"github.com/go-chi/chi/v5"
)

func HistoryTranslationRouter(router chi.Router) {
	router.Get("/generate", func(w http.ResponseWriter, r *http.Request) {
		aiService := aiservice.NewOpenAIAIService()
		historyTranslation := NewHistoryTranslation(aiService)

		story, err := historyTranslation.GenerateStory()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(story)
	})

	router.Post("/evaluate", func(w http.ResponseWriter, r *http.Request) {
		aiService := aiservice.NewOpenAIAIService()
		historyTranslation := NewHistoryTranslation(aiService)

		story := r.FormValue("story")
		userTranslation := r.FormValue("userTranslation")
		evaluation, err := historyTranslation.EvaluateTranslation(story, userTranslation)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(evaluation)
	})
}
