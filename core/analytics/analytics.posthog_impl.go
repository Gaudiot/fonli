package analytics

import (
	"gaudiot.com/fonli/core"
	"github.com/posthog/posthog-go"
)

const (
	posthogEndpoint = "https://us.i.posthog.com"
)

type posthogAnalyticsService struct {
	client *posthog.Client
}

func NewPosthogAnalyticsService() *posthogAnalyticsService {
	return &posthogAnalyticsService{}
}

func (ph *posthogAnalyticsService) Init() error {
	envConfig := core.GetEnvConfig()
	posthogApiKey := envConfig.PostHogAPIKey

	client, err := posthog.NewWithConfig(posthogApiKey, posthog.Config{Endpoint: posthogEndpoint})
	if err != nil {
		return err
	}

	Client = client
	return nil
}

func (s *posthogAnalyticsService) Close() error {
	if s.client != nil {
		err := (*s.client).Close()
		s.client = nil
		return err
	}
	return nil
}
