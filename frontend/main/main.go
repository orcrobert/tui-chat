package main

import (
	"log"

	"tui-chat/auth"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	initialMode := auth.NewModel()
	p := tea.NewProgram(initialMode)

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
