package main

import (
	"log"
	"net/http"
	"time"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	refreshtoken_repo "gaudiot.com/fonli/base/repositories/refresh_token"
	user_repo "gaudiot.com/fonli/base/repositories/user"
	"gaudiot.com/fonli/core"
	"gaudiot.com/fonli/core/analytics"
	"gaudiot.com/fonli/core/database"
	"gaudiot.com/fonli/core/middlewares"
	"gaudiot.com/fonli/core/security/password"
	"gaudiot.com/fonli/core/security/tokens"
	"gaudiot.com/fonli/src/authentication"
	"gaudiot.com/fonli/src/exercises"
	storytranslation "gaudiot.com/fonli/src/exercises/story_translation"
	wordconjugationexercise "gaudiot.com/fonli/src/exercises/word_conjugation"
	wordtranslationexercise "gaudiot.com/fonli/src/exercises/word_translation"
	"gaudiot.com/fonli/src/user_settings"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

const (
	maxBytes = 1 << 20 // 1MB

	requestTimeout = 30 * time.Second

	authRateLimit    = 10
	defaultRateLimit = 20
)

func main() {
	err := core.LoadEnvConfig()
	if err != nil {
		log.Fatal(err)
	}

	envConfig := core.GetEnvConfig()
	db, err := database.Connect(envConfig.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = analytics.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer analytics.Close()

	tokenService := tokens.NewTokenService([]byte(envConfig.JWTSecret))
	passwordService := &password.BCryptPasswordService{}
	aiService := &aiservice.OpenAIAIService{}
	var userRepository user_repo.UserRepository = user_repo.NewPgxUserRepository(db)
	var refreshTokenRepository refreshtoken_repo.RefreshTokenRepository = refreshtoken_repo.NewPgxRefreshTokenRepository(db)
	var userSettingsService *user_settings.UserSettingsService = user_settings.NewUserSettingsService(userRepository, aiService)
	var authService *authentication.AuthService = authentication.NewAuthService(tokenService, passwordService, userRepository, refreshTokenRepository)

	wordTranslation := wordtranslationexercise.NewWordTranslation(aiService)
	wordConjugation := wordconjugationexercise.NewWordConjugation(aiService)
	storyTranslation := storytranslation.NewStoryTranslation(aiService)

	router := chi.NewRouter()
	router.Use(middlewares.MaxBytesMiddleware(maxBytes))
	router.Use(middleware.Timeout(requestTimeout))

	authRateLimiter := httprate.LimitByIP(authRateLimit, time.Minute)
	defaultRateLimiter := httprate.LimitByIP(defaultRateLimit, time.Minute)

	router.Group(func(r chi.Router) {
		r.Use(authRateLimiter)
		r.Route("/auth", authentication.AuthenticationRouter(authService))
	})

	router.Group(func(r chi.Router) {
		r.Use(defaultRateLimiter)
		r.Route("/user", user_settings.UserSettingsRouter(userSettingsService, tokenService))

		r.Route("/exercises", exercises.ExercisesRouter(wordConjugation, wordTranslation, storyTranslation, tokenService))
	})

	log.Printf("Server is running on port :%s", envConfig.Port)
	http.ListenAndServe(":"+envConfig.Port, router)
}
