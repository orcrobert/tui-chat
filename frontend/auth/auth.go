package auth

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"tui-chat/chatview"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	usernameInput textinput.Model
	passwordInput textinput.Model
	width         int
	height        int
}

func NewModel() model {
	tiname := textinput.New()
	tiname.Placeholder = "Username"
	tiname.Focus()
	tiname.CharLimit = 10
	tiname.Width = 50

	tipass := textinput.New()
	tipass.Placeholder = "Password"
	tipass.CharLimit = 20
	tipass.Width = 50

	return model{
		usernameInput: tiname,
		passwordInput: tipass,
	}
}

func (m model) GetUsername() string {
	return m.usernameInput.Value()
}

func ValidateUser(username, password string) bool {
	executablePath, err := filepath.Abs("../../build/validate")
	if err != nil {
		fmt.Println("Error determining executable path:", err)
		return false
	}

	cmd := exec.Command(executablePath)
	cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%s\n%s\n", username, password))

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return string(output) == "Valid!\n"
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.usernameInput, cmd = m.usernameInput.Update(msg)
	m.passwordInput, cmd = m.passwordInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" || msg.String() == "tab" {
			if m.usernameInput.Focused() {
				m.usernameInput.Blur()
				m.passwordInput.Focus()

			} else if m.passwordInput.Focused() {
				if ValidateUser(m.usernameInput.Value(), m.passwordInput.Value()) {
					m.passwordInput.Blur()
					fmt.Printf("Authenticated as %s\n", m.usernameInput.Value())
					chatModel := chatview.NewModel(m.usernameInput.Value())

					chatProgram := tea.NewProgram(chatModel)

					if err := chatProgram.Start(); err != nil {
						log.Fatal(err)
					}

					return m, tea.Quit
				} else {
					fmt.Println("Invalid login credentials!")
					return m, tea.Quit
				}
			}
		} else if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.usernameInput.Width = m.width / 2
		m.passwordInput.Width = m.width / 2
	}

	return m, cmd
}

func (m model) View() string {
	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("black")).
		Bold(true).
		Padding(1, 1)

	centeredUsername := fmt.Sprintf("%s\n", m.usernameInput.View())
	centeredPassword := fmt.Sprintf("%s", m.passwordInput.View())

	return inputStyle.Render(centeredUsername + centeredPassword)
}
