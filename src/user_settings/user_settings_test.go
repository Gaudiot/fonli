package user_settings

import (
	"errors"
	"testing"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	user_repo "gaudiot.com/fonli/base/repositories/user"
)

func newTestUserSettingsService() (*UserSettingsService, *aiservice.AIServiceMock, *user_repo.UserRepositoryMock) {
	mockAI := &aiservice.AIServiceMock{}
	mockRepo := &user_repo.UserRepositoryMock{Users: map[string]*user_repo.User{
		"user1": {
			ID:       "user1",
			Email:    "user@example.com",
			Username: "testuser",
		},
	}}
	return NewUserSettingsService(mockRepo, mockAI), mockAI, mockRepo
}

// MARK: - UpdateUserLifestyle

func TestUpdateUserLifestyleSuccess(t *testing.T) {
	service, mockAI, _ := newTestUserSettingsService()

	mockAI.PromptFunc = func(prompt string) (string, error) {
		return "gym, soccer, nurse", nil
	}

	err := service.UpdateUserLifestyle("user1", "I like to go to the gym, play soccer and work as a nurse")
	if err != nil {
		t.Fatalf("UpdateUserLifestyle(); Unexpected error: %v", err)
	}
}

func TestUpdateUserLifestyleAIServiceError(t *testing.T) {
	service, mockAI, _ := newTestUserSettingsService()
	aiErr := errors.New("ai service unavailable")

	mockAI.PromptFunc = func(prompt string) (string, error) {
		return "", aiErr
	}

	err := service.UpdateUserLifestyle("user1", "some text")
	if err == nil {
		t.Fatalf("UpdateUserLifestyle(); Wanted error, got nil")
	}
	if !errors.Is(err, aiErr) {
		t.Errorf("UpdateUserLifestyle(); Wanted error %v, got %v", aiErr, err)
	}
}

// MARK: - GetUserLifestyle (TDD — not yet implemented)

func TestGetUserLifestyleSuccess(t *testing.T) {
	service, _, mockRepo := newTestUserSettingsService()
	mockRepo.Users["user1"].Lifestyle = "I like to go to the gym and play soccer"
	mockRepo.Users["user1"].LifestyleTopics = "gym, soccer"

	lifestyle, err := service.GetUserLifestyle("user1")
	if err != nil {
		t.Fatalf("GetUserLifestyle(); Unexpected error: %v", err)
	}
	if lifestyle == "" {
		t.Errorf("GetUserLifestyle(); Wanted non-empty lifestyle, got empty")
	}
}

func TestGetUserLifestyleUserNotFound(t *testing.T) {
	service, _, _ := newTestUserSettingsService()

	_, err := service.GetUserLifestyle("nonexistent")
	if err == nil {
		t.Fatalf("GetUserLifestyle(); Wanted error for nonexistent user, got nil")
	}
}
