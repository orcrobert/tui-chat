// ARE ERORI AAAAAAAAA
package chatview

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	conn      *exec.Cmd
	messages  []string
	textInput textinput.Model
	width     int
	height    int
	username  string
}

func NewModel(username string) (model, error) {
	// Start the client program (client.cpp) through exec
	cmd := exec.Command("./client", username) // Make sure `client.cpp` is compiled to `client` and in the same directory

	// Start the client process
	err := cmd.Start()
	if err != nil {
		return model{}, fmt.Errorf("failed to start client: %w", err)
	}

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 50

	return model{
		username:  username,
		conn:      cmd,
		messages:  []string{},
		textInput: ti,
	}, nil
}

func (m model) Init() tea.Cmd {
	// Start listening for messages from the client (C++ process)
	return tea.Batch(textinput.Blink, listenForMessages(m.conn))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.textInput.Value() != "" {
				m.messages = append(m.messages, fmt.Sprintf("%s: %s", m.username, m.textInput.Value()))
				message := m.textInput.Value()

				// Send message to the C++ client (client.cpp)
				fmt.Fprintf(m.conn.Stdin, "%s\n", message)
				m.textInput.SetValue("")
			}
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	case string:
		// Append incoming messages to the list
		m.messages = append(m.messages, msg)
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

// Listen for incoming messages from the C++ client
func listenForMessages(cmd *exec.Cmd) tea.Cmd {
	return func() tea.Msg {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}

		reader := bufio.NewReader(stdout)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Sprintf("Error reading message: %v", err)
			}

			message = strings.TrimSpace(message)
			// Return the message to be appended to the UI
			return message
		}
	}
}
