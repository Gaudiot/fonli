package main

import (
	"log"
	"net/http"

	"gaudiot.com/fonli/core"
	wordtranslationexercise "gaudiot.com/fonli/src/word_translation"
	"github.com/go-chi/chi/v5"
)

func main() {
	err := core.LoadEnvConfig()
	if err != nil {
		log.Fatal(err)
	}

	envConfig := core.GetEnvConfig()

	router := chi.NewRouter()
	router.Route("/word-translation", wordtranslationexercise.WordTranslationRouter)

	log.Printf("Server is running on port :%s", envConfig.Port)
	http.ListenAndServe(":"+envConfig.Port, router)
}
