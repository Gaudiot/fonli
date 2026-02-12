package core

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	OpenAIKey string
	// DB    PostgresConfig
	// Port  string
}

// type LogConfig struct {
// 	Style string
// 	Level string
// }

// type PostgresConfig struct {
// 	Username string
// 	Password string
// 	URL      string
// 	Port     string
// }

func getEnvValue(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", errors.New(key + " environment variable is required and cannot be empty")
	}
	return value, nil
}

func LoadEnvConfig() (*EnvConfig, error) {
	godotenv.Load("../.env")

	openAIKey, err := getEnvValue("OPENAI_API_KEY")
	if err != nil {
		return nil, err
	}

	config := &EnvConfig{
		OpenAIKey: openAIKey,
	}

	return config, nil
}
