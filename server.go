package main

import (
	"log"
	"net/http"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	refreshtoken_repo "gaudiot.com/fonli/base/repositories/refresh_token"
	user_repo "gaudiot.com/fonli/base/repositories/user"
	"gaudiot.com/fonli/core"
	"gaudiot.com/fonli/core/security/password"
	"gaudiot.com/fonli/core/security/tokens"
	"gaudiot.com/fonli/src/authentication"
	storytranslation "gaudiot.com/fonli/src/story_translation"
	"gaudiot.com/fonli/src/user_settings"
	wordconjugationexercise "gaudiot.com/fonli/src/word_conjugation"
	wordtranslationexercise "gaudiot.com/fonli/src/word_translation"
	"github.com/go-chi/chi/v5"
)

func main() {
	err := core.LoadEnvConfig()
	if err != nil {
		log.Fatal(err)
	}

	envConfig := core.GetEnvConfig()

	// TODO: replace mocks with real repository implementations
	tokenService := tokens.NewTokenService([]byte(envConfig.JWTSecret))
	passwordService := &password.BCryptPasswordService{}
	userRepository := &user_repo.UserRepositoryMock{Users: make(map[string]*user_repo.User)}
	aiService := &aiservice.OpenAIAIService{}
	userSettingsService := user_settings.NewUserSettingsService(userRepository, aiService)
	refreshTokenRepository := &refreshtoken_repo.RefreshTokenRepositoryMock{RefreshTokens: make(map[string]*refreshtoken_repo.RefreshToken)}
	authService := authentication.NewAuthService(*tokenService, passwordService, userRepository, refreshTokenRepository)

	router := chi.NewRouter()
	router.Route("/auth", authentication.AuthenticationRouter(authService))

	router.Route("/user", user_settings.UserSettingsRouter(userSettingsService, tokenService))

	router.Route("/word-translation", wordtranslationexercise.WordTranslationRouter)
	router.Route("/word-conjugation", wordconjugationexercise.WordConjugationRouter)
	router.Route("/history-translation", storytranslation.StoryTranslationRouter)

	log.Printf("Server is running on port :%s", envConfig.Port)
	http.ListenAndServe(":"+envConfig.Port, router)
}
