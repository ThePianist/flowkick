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

type gotProjectSuccessMsg []projectNames

type projectNames struct {
	Name string `json:"name"`
}

func getProjects(store *store.Store) ([]string, error) {
	return store.GetScopes()
}

type ProjectSearchModel struct {
	textInput textinput.Model
	help      help.Model
	keymap    projectKeymap
}

type projectKeymap struct{}

func (k projectKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		key.NewBinding(key.WithKeys("↑", "up"), key.WithHelp("↑", "next")),
		key.NewBinding(key.WithKeys("↓", "down"), key.WithHelp("↓", "prev")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}

func (k projectKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func NewProjectSearchModel(entry, selectedType string, store *store.Store) ProjectSearchModel {
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
	km := projectKeymap{}

	// Fetch dynamic suggestions from DB
	var suggestions []string
	if store != nil {
		suggestions, _ = getProjects(store)
	}
	ti.SetSuggestions(suggestions)

	return ProjectSearchModel{
		textInput: ti,
		help:      h,
		keymap:    km,
	}
}

func (m ProjectSearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ProjectSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case gotProjectSuccessMsg:
		var suggestions []string
		for _, r := range msg {
			suggestions = append(suggestions, r.Name)
		}
		m.textInput.SetSuggestions(suggestions)
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m ProjectSearchModel) View() string {
	var (
		textInputTitleStyle = lipgloss.NewStyle().MarginTop(1).MarginLeft(2).Bold(true).Foreground(lipgloss.Color("63"))
		textInputInputStyle = lipgloss.NewStyle().PaddingLeft(2)
	)

	return fmt.Sprintf(
		textInputTitleStyle.Render("📌 Project: %s\n\n%s\n\n"),
		textInputInputStyle.Render(m.textInput.View()),
		m.help.View(m.keymap),
	)
}
