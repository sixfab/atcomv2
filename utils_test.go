package atcom

import (
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {

	cmd := "echo"
	arg := "hello"
	expectedOutput := "hello\n"

	output, err := RunCommand(cmd, arg)
	if err != nil {
		t.Errorf("RunCommand returned an error: %v", err)
	}
	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}
}
