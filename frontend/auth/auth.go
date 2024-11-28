package auth

import (
	"fmt"

	"time"

	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

type model struct {
	usernameInput textinput.Model
	passwordInput textinput.Model
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
				m.passwordInput.Blur()
				time.Sleep(2 * time.Second)
				return m, tea.Quit
			}
		} else if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
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
