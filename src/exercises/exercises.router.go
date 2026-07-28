package exercises

import (
	"net/http"

	"gaudiot.com/fonli/core/middlewares"
	"gaudiot.com/fonli/core/security/tokens"
	storytranslation "gaudiot.com/fonli/src/exercises/story_translation"
	verbconjugationexercise "gaudiot.com/fonli/src/exercises/verb_conjugation"
	vocabularyexercise "gaudiot.com/fonli/src/exercises/vocabulary"
	"github.com/go-chi/chi/v5"
)

// MARK: - Router

func ExercisesRouter(
	wc *verbconjugationexercise.WordConjugation,
	wt *vocabularyexercise.VocabularyExercise,
	st *storytranslation.StoryTranslation,
	ts tokens.TokenService,
) func(chi.Router) {
	return func(router chi.Router) {
		router.Use(middlewares.AuthMiddleware(ts))

		router.Route("/verb-conjugation", verbconjugationexercise.WordConjugationRouter(wc))
		router.Route("/vocabulary", vocabularyexercise.VocabularyRouter(wt))
		router.Route("/story-translation", storytranslation.StoryTranslationRouter(st))
		router.Get("/", handleRoot())
	}
}

// Optional: root handler just to confirm this router is alive
func handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true,"message":"Fonli Exercises API"}`))
	}
}
