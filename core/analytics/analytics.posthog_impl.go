package analytics

import (
	"errors"

	"gaudiot.com/fonli/core"
	"github.com/posthog/posthog-go"
)

const (
	posthogEndpoint = "https://us.i.posthog.com"
)

type posthogAnalyticsService struct {
	client posthog.Client
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

	ph.client = client
	return nil
}

func (ph *posthogAnalyticsService) Close() error {
	if ph.client != nil {
		err := ph.client.Close()
		ph.client = nil
		return err
	}
	return nil
}

func (ph *posthogAnalyticsService) Register(eventID string, properties map[string]any) error {
	if ph.client == nil {
		return errors.New("posthog client is not initialized")
	}

	// É necessário passar um UserId para o evento no PostHog. Se não houver, define um padrão.
	distinctId, ok := properties["distinct_id"].(string)
	if !ok || distinctId == "" {
		distinctId = "anonymous"
	}

	return ph.client.Enqueue(&posthog.Capture{
		Event:      eventID,    // EventID
		DistinctId: distinctId, // UserId
		Properties: properties,
	})
}
