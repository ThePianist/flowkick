package cmd

import (
	"fmt"

	"github.com/ThePianist/flowkick/store"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getIssues(store *store.Store) ([]string, error) {
	return store.GetTickets()
}

type IssueSearchModel struct {
	textInput   textinput.Model
	help        help.Model
	keymap      issueKeymap
	store       *store.Store
	suggestions []string
}

type issueKeymap struct{}

func (k issueKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		key.NewBinding(key.WithKeys("↑", "up"), key.WithHelp("↑", "next")),
		key.NewBinding(key.WithKeys("↓", "down"), key.WithHelp("↓", "prev")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}

func (k issueKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func NewIssueSearchModel(entry, selectedType string, store *store.Store) IssueSearchModel {
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
	km := issueKeymap{}

	// Fetch dynamic suggestions from DB
	var suggestions []string
	if store != nil {
		suggestions, _ = getIssues(store)
	}
	ti.SetSuggestions(suggestions)

	return IssueSearchModel{
		textInput:   ti,
		help:        h,
		keymap:      km,
		store:       store,
		suggestions: suggestions,
	}
}

func (m IssueSearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m IssueSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			// Save new issue/ticket if not already in DB
			input := m.textInput.Value()
			if input != "" && m.suggestions != nil {
				found := false
				for _, s := range m.suggestions {
					if s == input {
						found = true
						break
					}
				}
				if !found {
					// Save new ticket
					if m.store != nil {
						_ = m.store.SaveTicket(0, input, 0) // entryID can be set if needed
					}
				}
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m IssueSearchModel) View() string {
	var (
		textInputTitleStyle = lipgloss.NewStyle().MarginTop(1).MarginLeft(2).Bold(true).Foreground(lipgloss.Color("63"))
		textInputInputStyle = lipgloss.NewStyle().PaddingLeft(2)
	)

	return fmt.Sprintf(
		textInputTitleStyle.Render("📌 Related issue: %s\n\n%s\n\n"),
		textInputInputStyle.Render(m.textInput.View()),
		m.help.View(m.keymap),
	)
}
