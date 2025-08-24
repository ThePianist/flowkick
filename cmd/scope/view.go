package scope

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

func (m ScopeModel) View() string {
	var (
		textInputTitleStyle = lipgloss.NewStyle().MarginTop(1).MarginLeft(2).Bold(true).Foreground(lipgloss.Color("63"))
		textInputInputStyle = lipgloss.NewStyle().PaddingLeft(2)
	)

	return fmt.Sprintf(
		textInputTitleStyle.Render("📌 Project: %s\n\n%s\n\n"),
		textInputInputStyle.Render(m.Input.View()),
		m.Help.View(m.Keymap),
	)
}

func (k ScopeKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		key.NewBinding(key.WithKeys("↑", "up"), key.WithHelp("↑", "next")),
		key.NewBinding(key.WithKeys("↓", "down"), key.WithHelp("↓", "prev")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}

func (k ScopeKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
