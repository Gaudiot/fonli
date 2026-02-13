package wordtranslationexercise

import (
	"encoding/json"
	"net/http"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	"github.com/go-chi/chi/v5"
)

func WordTranslationRouter(router chi.Router) {
	router.Get("/native-to-foreign", func(w http.ResponseWriter, r *http.Request) {
		aiService := aiservice.NewOpenAIAIService()
		wordTranslation := NewWordTranslation(aiService)

		exercises, err := wordTranslation.NativeToForeignExercise(10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(exercises)
	})

	router.Get("/foreign-to-native", func(w http.ResponseWriter, r *http.Request) {
		aiService := aiservice.NewOpenAIAIService()
		wordTranslation := NewWordTranslation(aiService)

		exercises, err := wordTranslation.ForeignToNativeExercise(10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(exercises)
	})
}
