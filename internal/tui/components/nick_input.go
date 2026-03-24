package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NickInput struct {
	input     textinput.Model
	done      bool
	width     int
	height    int
	theme     lipgloss.Color
	onConfirm func(string)
}

func NewNickInput(theme lipgloss.Color) *NickInput {
	ti := textinput.New()
	ti.Placeholder = "Enter your nickname"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 30

	return &NickInput{
		input: ti,
		theme: theme,
	}
}

func (n *NickInput) Update(msg tea.Msg) (*NickInput, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			nick := strings.TrimSpace(n.input.Value())
			if nick != "" && n.onConfirm != nil {
				n.onConfirm(nick)
				n.done = true
			}
			return n, nil
		case "esc":
			n.done = true
			return n, nil
		}
	}

	n.input, cmd = n.input.Update(msg)
	return n, cmd
}

func (n *NickInput) View() string {
	if n.done {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(n.theme).
		Bold(true).
		MarginBottom(1)

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00"))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("WELCOME TO T-MESS"),
		lipgloss.NewStyle().MarginBottom(1).Render("Please set your nickname:"),
		promptStyle.Render("> "),
		n.input.View(),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			MarginTop(1).
			Render("Press Enter to continue"),
	)

	return lipgloss.Place(
		n.width, n.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (n *NickInput) SetSize(width, height int) {
	n.width = width
	n.height = height
}

func (n *NickInput) SetOnConfirm(f func(string)) {
	n.onConfirm = f
}

func (n *NickInput) IsDone() bool {
	return n.done
}
