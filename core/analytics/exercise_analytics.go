package analytics

const EventExerciseInvoked = "exercise_invoked"

type ExerciseType string

const (
	ExerciseWordTranslationNativeToForeign ExerciseType = "word_translation_native_to_foreign"
	ExerciseWordTranslationForeignToNative ExerciseType = "word_translation_foreign_to_native"
	ExerciseWordConjugation                ExerciseType = "word_conjugation"
	ExerciseStoryTranslationGenerate       ExerciseType = "story_translation_generate"
	ExerciseStoryTranslationEvaluate       ExerciseType = "story_translation_evaluate"
)

type ExerciseOutcome string

const (
	ExerciseOutcomeSuccess         ExerciseOutcome = "success"
	ExerciseOutcomeValidationError ExerciseOutcome = "validation_error"
	ExerciseOutcomeInternalError   ExerciseOutcome = "internal_error"
)

func TrackExerciseInvocation(userID string, exerciseType ExerciseType, outcome ExerciseOutcome, err ...error) {
	if Client == nil {
		return
	}
	if userID == "" {
		return
	}
	props := map[string]interface{}{
		"distinct_id":     userID,
		"exercise_type":   string(exerciseType),
		"outcome":         string(outcome),
		"success":         outcome == ExerciseOutcomeSuccess,
		"failed_internal": outcome == ExerciseOutcomeInternalError,
	}
	if len(err) > 0 && err[0] != nil {
		props["error_message"] = err[0].Error()
	}
	_ = Client.Register(EventExerciseInvoked, props)
}
