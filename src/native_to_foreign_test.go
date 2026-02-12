package src

import "testing"

func TestCreateExercise(t *testing.T) {
	exercisesQuantity := 10

	got, err := CreateNativeToForeignExercise(exercisesQuantity)

	if err != nil {
		t.Errorf("CreateExercise(%d) should not return an error, but got %v", exercisesQuantity, err)
	}

	if len(got.Questions) != exercisesQuantity {
		t.Errorf("CreateExercise(%d) = %s,  should have %d questions, but has %d", exercisesQuantity, got, exercisesQuantity, len(got.Questions))
	}
}
