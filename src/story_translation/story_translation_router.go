package storytranslation

import (
	"encoding/json"
	"net/http"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	"github.com/go-chi/chi/v5"
)

func StoryTranslationRouter(router chi.Router) {
	router.Get("/generate", func(w http.ResponseWriter, r *http.Request) {
		nativeLanguageCode := r.URL.Query().Get("nl")
		foreignLanguageCode := r.URL.Query().Get("fl")

		aiService := aiservice.NewOpenAIAIService()
		storyTranslation := NewStoryTranslation(aiService)

		story, err := storyTranslation.GenerateStory(nativeLanguageCode, foreignLanguageCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(story)
	})

	router.Post("/evaluate", func(w http.ResponseWriter, r *http.Request) {
		nativeLanguageCode := r.URL.Query().Get("nl")
		foreignLanguageCode := r.URL.Query().Get("fl")

		aiService := aiservice.NewOpenAIAIService()
		storyTranslation := NewStoryTranslation(aiService)

		story := r.FormValue("story")
		userTranslation := r.FormValue("userTranslation")
		evaluation, err := storyTranslation.EvaluateTranslation(story, userTranslation, nativeLanguageCode, foreignLanguageCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(evaluation)
	})
}
