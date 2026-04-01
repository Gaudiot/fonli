package user_settings

import (
	"errors"
	"fmt"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	user_repository "gaudiot.com/fonli/base/repositories/user"
)

type UserSettingsService struct {
	userRepository user_repository.UserRepository
	aiService      aiservice.AIService
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func NewUserSettingsService(userRepository user_repository.UserRepository, aiService aiservice.AIService) *UserSettingsService {
	return &UserSettingsService{
		userRepository: userRepository,
		aiService:      aiService,
	}
}

func (us *UserSettingsService) GetUserLifestyle(userID string) (string, error) {
	user, err := us.userRepository.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotFound
	}
	return user.LifestyleTopics, nil
}

func (us *UserSettingsService) UpdateUserLifestyle(userID, text string) error {
	prompt := fmt.Sprintf(
		`The user is describing his daily lifestyle. This can include the activities he does, the places he goes, the people he spends time with, the things he likes to do, etc.
		Based on the text, I need you to return me lifestyle topics.
		Lifestyle topics is a string in CSV format with each value being a user's different lifestyle topic.

		As an example, if the user is describing his daily lifestyle as: "I like to go to the gym, play soccer and work as a nurse"
		The lifestyle topics would be: "gym, soccer, nurse"

		Return only the lifestyle points, no other text.
		The text is: %s
		`,
		text,
	)

	response, err := us.aiService.Prompt(prompt)
	if err != nil {
		return err
	}

	us.userRepository.UpdateUserLifestyle(userID, text, response)

	return nil
}
