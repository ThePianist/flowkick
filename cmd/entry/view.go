package entry

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	textInputTitleStyle = lipgloss.NewStyle().MarginTop(1).MarginLeft(2).Bold(true).Foreground(lipgloss.Color("63"))
	textInputInputStyle = lipgloss.NewStyle().PaddingLeft(2)
	textInputQuitStyle  = lipgloss.NewStyle().MarginLeft(2)
)

func (m EntryModel) View() string {
	title := textInputTitleStyle.Render("💭 What's on your mind?")
	input := textInputInputStyle.Render(m.Input.View())
	quit := textInputQuitStyle.Render("(esc to quit)")

	return title + "\n\n" + input + "\n\n" + quit
}
