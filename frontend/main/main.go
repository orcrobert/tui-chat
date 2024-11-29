package main

import (
	"fmt"
	"log"
	"tui-chat/auth"
	"tui-chat/chatview"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	authModel := auth.NewModel()
	p := tea.NewProgram(authModel)

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	username := "aaaa"
	if username != "" {
		fmt.Printf("Authenticated as %s\n", username)
		chatModel := chatview.NewModel(username)

		chatProgram := tea.NewProgram(chatModel)

		if err := chatProgram.Start(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Authentication failed or returned unexpected result")
	}
}
