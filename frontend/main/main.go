package main

import (
	"log"
	"tui-chat/auth"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	authModel := auth.NewModel()
	p := tea.NewProgram(authModel)

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
