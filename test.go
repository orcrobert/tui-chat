package main
import (
	"fmt"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type model struct {
	messages  []string
	textInput textinput.Model
	width     int
	height    int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 50

	return model{
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
			m.messages = append(m.messages, m.textInput.Value())
			m.textInput.SetValue("")
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

    // Calculate the height for messages, reserving 3 lines for the input box
    messageAreaHeight := m.height - 3 // 1 for input box, 2 for padding (borders or separators)
    if messageAreaHeight < 0 {
        messageAreaHeight = 0 // Prevent negative heights
    }

    // Calculate the starting index for the messages
    start := 0
    if len(m.messages) > messageAreaHeight {
        start = len(m.messages) - messageAreaHeight
    }

    // Append messages to the buffer
    visibleMessages := m.messages[start:]
    for _, msg := range visibleMessages {
        b.WriteString(fmt.Sprintf("%s\n", msg))
    }

    // Fill any empty space with blank lines to push the input box to the bottom
    blankLines := messageAreaHeight - len(visibleMessages)
    for i := 0; i < blankLines; i++ {
        b.WriteString("\n")
    }

    // Add separator line and input box at the bottom
    inputLine := strings.Repeat("─", m.width)
    b.WriteString(inputLine + "\n")
    b.WriteString(m.textInput.View())

    return b.String()
}



func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

