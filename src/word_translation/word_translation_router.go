package wordtranslationexercise

import (
	"encoding/json"
	"net/http"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	"github.com/go-chi/chi/v5"
)

func WordTranslationRouter(router chi.Router) {
	router.Get("/native-to-foreign", func(w http.ResponseWriter, r *http.Request) {
		nl := r.URL.Query().Get("nl")
		fl := r.URL.Query().Get("fl")

		aiService := aiservice.NewOpenAIAIService()
		wordTranslation := NewWordTranslation(aiService)

		exercises, err := wordTranslation.NativeToForeignExercise(10, nl, fl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercises)
	})

	router.Get("/foreign-to-native", func(w http.ResponseWriter, r *http.Request) {
		nl := r.URL.Query().Get("nl")
		fl := r.URL.Query().Get("fl")

		aiService := aiservice.NewOpenAIAIService()
		wordTranslation := NewWordTranslation(aiService)

		exercises, err := wordTranslation.ForeignToNativeExercise(10, fl, nl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercises)
	})
}
