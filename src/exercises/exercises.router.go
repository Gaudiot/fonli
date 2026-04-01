package exercises

import (
	"net/http"

	"gaudiot.com/fonli/core/middlewares"
	"gaudiot.com/fonli/core/security/tokens"
	storytranslation "gaudiot.com/fonli/src/story_translation"
	wordconjugationexercise "gaudiot.com/fonli/src/word_conjugation"
	wordtranslationexercise "gaudiot.com/fonli/src/word_translation"
	"github.com/go-chi/chi/v5"
)

// MARK: - Router

func ExercisesRouter(
	wc *wordconjugationexercise.WordConjugation,
	wt *wordtranslationexercise.WordTranslation,
	st *storytranslation.StoryTranslation,
	ts *tokens.TokenService,
) func(chi.Router) {
	return func(router chi.Router) {
		router.Use(middlewares.AuthMiddleware(ts))

		router.Route("/word-conjugation", wordconjugationexercise.WordConjugationRouter(wc, ts))
		router.Route("/word-translation", wordtranslationexercise.WordTranslationRouter(wt, ts))
		router.Route("/story-translation", storytranslation.StoryTranslationRouter(st, ts))
		router.Get("/", handleRoot())
	}
}

// Optional: root handler just to confirm this router is alive
func handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true,"message":"Fonli Exercises API"}`))
	}
}
