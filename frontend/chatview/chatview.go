package chatview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	messages  []string
	textInput textinput.Model
	width     int
	height    int
	username  string
}

func NewModel(username string) model {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 50
	return model{
		username:  username,
		messages:  []string{},
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.textInput.Value() != "" {
				m.messages = append(m.messages, fmt.Sprintf("%s: %s", m.username, m.textInput.Value()))
				m.textInput.SetValue("")
			}
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	messageAreaHeight := m.height - 3 
	if messageAreaHeight < 0 {
		messageAreaHeight = 0 
	}

	start := 0
	if len(m.messages) > messageAreaHeight {
		start = len(m.messages) - messageAreaHeight
	}

	visibleMessages := m.messages[start:]
	for _, msg := range visibleMessages {
		b.WriteString(fmt.Sprintf("%s\n", msg))
	}

	blankLines := messageAreaHeight - len(visibleMessages)
	for i := 0; i < blankLines; i++ {
		b.WriteString("\n")
	}

	inputLine := strings.Repeat("─", m.width)
	b.WriteString(inputLine + "\n")
	b.WriteString(m.textInput.View())

	return b.String()
}
