package scope

import (
	"github.com/ThePianist/flowkick/store"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type ScopeKeymap struct{}

type ScopeModel struct {
	Input        textinput.Model
	Help         help.Model
	Keymap       ScopeKeymap
	ErrorMessage string // User-facing error message
}

func NewScopeModel(entry, selectedType string, store *store.Store) ScopeModel {
	ti := textinput.New()
	ti.Placeholder = "(press ↵ to skip)"
	ti.Prompt = ""
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 60
	ti.ShowSuggestions = true

	h := help.New()
	km := ScopeKeymap{}

	// Fetch dynamic suggestions from DB
	var suggestions []string
	if store != nil {
		scopes, _ := getScopes(store)
		for _, scope := range scopes {
			suggestions = append(suggestions, scope.Name)
		}
	}
	ti.SetSuggestions(suggestions)

	return ScopeModel{
		Input:  ti,
		Help:   h,
		Keymap: km,
	}
}

func getScopes(store *store.Store) ([]store.Scope, error) {
	return store.GetScopes()
}

func (m ScopeModel) GetValue() string {
	return m.Input.Value()
}
