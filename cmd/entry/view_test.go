package entry

import (
	"errors"
	"strings"
	"testing"
)

func TestEntryModel_View(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{
			name:     "empty input",
			input:    "",
			hasError: false,
		},
		{
			name:     "with input text",
			input:    "Fixed weird cache bug after 2hrs",
			hasError: false,
		},
		{
			name:     "with error",
			input:    "test input",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a proper EntryModel
			m := NewEntryModel()
			m.Input.SetValue(tt.input)

			if tt.hasError {
				m.Err = errors.New("some error")
			}

			got := m.View()

			// Test that the view contains expected elements
			if !strings.Contains(got, "💭 What's on your mind?") {
				t.Errorf("View() should contain title, got %v", got)
			}

			if !strings.Contains(got, "(esc to quit)") {
				t.Errorf("View() should contain quit instruction, got %v", got)
			}

			// For non-empty inputs, check that the input text appears
			if tt.input != "" && !strings.Contains(got, tt.input) {
				t.Errorf("View() should contain input text %q, got %v", tt.input, got)
			}
		})
	}
}
