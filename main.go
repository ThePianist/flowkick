package main

import (
	"log"

	"github.com/ThePianist/flowkick/cmd"
	"github.com/ThePianist/flowkick/logger"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize logging to a file
	logger.Init("flowkick.log")
	p := tea.NewProgram(cmd.InitialAppModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
