package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type gotIssueSuccessMsg []issueIds

type issueIds struct {
	Name string `json:"name"`
}

func getRepos() tea.Msg {
	// Add default Jira ticket issue names
	issueIds := []issueIds{
		{Name: "WBS-100"},
		{Name: "WBS-104"},
		{Name: "WBS-110"},
		{Name: "WBS-111"},
	}

	return gotIssueSuccessMsg(issueIds)
}

type IssueSearchModel struct {
	textInput textinput.Model
	help      help.Model
	keymap    issueKeymap
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

func NewIssueSearchModel(entry, selectedType string) IssueSearchModel {
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

	return IssueSearchModel{
		textInput: ti,
		help:      h,
		keymap:    km,
	}
}

func (m IssueSearchModel) Init() tea.Cmd {
	return tea.Batch(getRepos, textinput.Blink)
}

func (m IssueSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case gotIssueSuccessMsg:
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
