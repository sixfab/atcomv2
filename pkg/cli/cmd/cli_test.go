package cmd

import (
	"testing"
)

func TestExecute(t *testing.T) {
	t.Run("should run Execute", func(t *testing.T) {
		Execute()
	})
}
