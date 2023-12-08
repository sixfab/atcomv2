package atcom

import (
	"errors"
	"testing"
)

func TestGetMeaningfulPart(t *testing.T) {

	t.Run("should return error", func(t *testing.T) {
		_, _, result := GetMeaningfulPart(nil, "", "")

		expectedResult := errors.New("no ok response")

		if !errors.Is(expectedResult, result) {
			t.Errorf("Expected %s, got %s", expectedResult, result)
		}
	})
}
