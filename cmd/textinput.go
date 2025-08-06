package cmd

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

type Model struct {
	AddEntry textinput.Model
	Err      error
}

func InitialModel() Model {
	entryInput := textinput.New()
	entryInput.Placeholder = "Fixed weird cache bug after 2hrs"
	entryInput.Focus()
	entryInput.CharLimit = 0
	entryInput.Width = 156

	return Model{
		AddEntry: entryInput,
		Err:      nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.Err = msg
		return m, nil
	}

	m.AddEntry, cmd = m.AddEntry.Update(msg)
	return m, cmd
}

var (
	textInputTitleStyle = lipgloss.NewStyle().MarginTop(1).MarginLeft(2).Bold(true).Foreground(lipgloss.Color("63"))
	textInputInputStyle = lipgloss.NewStyle().PaddingLeft(2)
	textInputQuitStyle  = lipgloss.NewStyle().MarginLeft(2)
)

func (m Model) View() string {
	title := textInputTitleStyle.Render("💭 What's on your mind?")
	input := textInputInputStyle.Render(m.AddEntry.View())
	quit := textInputQuitStyle.Render("(esc to quit)")

	return title + "\n\n" + input + "\n\n" + quit
}
