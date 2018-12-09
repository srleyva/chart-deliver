package helpers_test

import (
	"fmt"
	"testing"

	. "github.com/srleyva/chart-deliver/pkg/helpers"
)

func TestRun(t *testing.T) {
	handler := NewHelmHandler()
	input := "hello"
	out, err := handler.Run("echo", input)
	if err != nil {
		t.Errorf("err returned where not expected: %s", err)
	}

	if string(out) != fmt.Sprintf("%s\n", input) { // Output returns newline
		t.Errorf("command called incorrectly:\n Actual: %s Expected: %s", out, input)
	}
}
