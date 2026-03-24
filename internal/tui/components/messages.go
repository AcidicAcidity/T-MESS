package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AcidicAcidity/t-mess/internal/messages"
)

type MessageView struct {
	viewport viewport.Model
	messages []*messages.Message
	width    int
	height   int
	theme    lipgloss.Color
}

func NewMessageView(theme lipgloss.Color) *MessageView {
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle()

	return &MessageView{
		viewport: vp,
		theme:    theme,
	}
}

func (m *MessageView) SetMessages(msgs []*messages.Message) {
	m.messages = msgs
	m.updateContent()
	m.viewport.GotoBottom()
}

func (m *MessageView) AddMessage(msg *messages.Message) {
	m.messages = append(m.messages, msg)
	m.updateContent()
	m.viewport.GotoBottom()
}

func (m *MessageView) updateContent() {
	var content strings.Builder

	for _, msg := range m.messages {
		content.WriteString(m.renderMessage(msg))
		content.WriteString("\n")
	}

	m.viewport.SetContent(content.String())
}

func (m *MessageView) renderMessage(msg *messages.Message) string {
	timeStr := msg.Timestamp.Format("15:04:05")

	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Width(8).
		Align(lipgloss.Right)

	nameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ffff")).
		Bold(true)

	if msg.IsOwn {
		nameStyle = nameStyle.Foreground(lipgloss.Color("#ff66cc"))
	}

	sender := msg.SenderID
	if sender == "system" {
		sender = "📢 System"
		nameStyle = nameStyle.Foreground(lipgloss.Color("#ffff00"))
	} else if len(sender) > 16 {
		sender = sender[:14] + ".."
	}

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	line := lipgloss.JoinHorizontal(
		lipgloss.Top,
		timeStyle.Render(timeStr),
		" ",
		nameStyle.Render(sender),
		": ",
		textStyle.Render(msg.Text),
	)

	return line
}

func (m *MessageView) Update(msg tea.Msg) (*MessageView, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *MessageView) View() string {
	if len(m.messages) == 0 {
		return lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("#666666")).
			Width(m.width).
			Height(m.height).
			Render("✨ No messages yet. Send a message to start the conversation. ✨")
	}
	return m.viewport.View()
}

func (m *MessageView) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height
	m.updateContent()
}

func (m *MessageView) Focus() {}

func (m *MessageView) Blur() {}
