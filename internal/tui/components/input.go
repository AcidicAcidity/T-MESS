package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputField struct {
	textinput textinput.Model
	width     int
	onSend    func(string)
}

func NewInputField() *InputField {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()
	ti.CharLimit = 4096
	ti.Prompt = ""
	ti.Width = 80

	return &InputField{
		textinput: ti,
	}
}

func (i *InputField) Update(msg tea.Msg) (*InputField, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			text := i.textinput.Value()
			if text != "" && i.onSend != nil {
				i.onSend(text)
				i.textinput.Reset()
			}
			return i, nil
		}
	}

	i.textinput, cmd = i.textinput.Update(msg)
	return i, cmd
}

func (i *InputField) View() string {
	inputLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render(">_ "),
		i.textinput.View(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().
			Width(i.width).
			Padding(1, 1).
			BorderTop(true).
			BorderTopForeground(lipgloss.Color("#00aa00")).
			Render(inputLine),
		lipgloss.NewStyle().
			Width(i.width).
			Padding(0, 1).
			Foreground(lipgloss.Color("#888888")).
			Render("  [Enter] send  |  [Ctrl+C] quit"),
	)
}

func (i *InputField) SetSize(width int) {
	i.width = width
	i.textinput.Width = width - 10
}

func (i *InputField) SetOnSend(f func(string)) {
	i.onSend = f
}

func (i *InputField) Focus() {
	i.textinput.Focus()
}

func (i *InputField) Blur() {
	i.textinput.Blur()
}
