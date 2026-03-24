package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AcidicAcidity/t-mess/internal/messages"
)

type RightPanel struct {
	chat             *messages.Chat
	connectionStatus string
	width            int
	height           int
	theme            lipgloss.Color
}

func NewRightPanel(theme lipgloss.Color) *RightPanel {
	return &RightPanel{
		connectionStatus: "● CONNECTING...",
		theme:            theme,
	}
}

func (r *RightPanel) SetChat(chat *messages.Chat) {
	r.chat = chat
}

func (r *RightPanel) SetConnectionStatus(status string) {
	r.connectionStatus = status
}

func (r *RightPanel) Update(msg tea.Msg) (*RightPanel, tea.Cmd) {
	return r, nil
}

func (r *RightPanel) View() string {
	if r.chat == nil {
		return ""
	}

	// Инфо о чате
	nameStyle := lipgloss.NewStyle().
		Foreground(r.theme).
		Bold(true).
		MarginBottom(1)

	// Статус подключения
	statusColor := lipgloss.Color("#ff0000")
	if strings.Contains(r.connectionStatus, "ONLINE") {
		statusColor = lipgloss.Color("#00ff00")
	} else if strings.Contains(r.connectionStatus, "CONNECTING") {
		statusColor = lipgloss.Color("#ffff00")
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(statusColor).
		MarginBottom(2)

	// Заглушка для эмодзи
	emojiTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(2).
		MarginBottom(1).
		Render("QUICK EMOJIS")

	emojis := "😊 😂 🔥 🚀 💯 ❤️ 🎉 🤔"
	emojiStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffff00"))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		nameStyle.Render(fmt.Sprintf("📌 %s", r.chat.Name)),
		statusStyle.Render(r.connectionStatus),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("─── CHAT INFO ───"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(fmt.Sprintf("Type: %s", r.chat.Type)),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Encryption: 🔒 E2EE (mock)"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Peers: 0 (local only)"),
		emojiTitle,
		emojiStyle.Render(emojis),
	)

	return lipgloss.NewStyle().
		Padding(1, 2).
		BorderLeft(true).
		BorderLeftForeground(lipgloss.Color("#00aa00")).
		Width(r.width).
		Render(content)
}

func (r *RightPanel) SetSize(width, height int) {
	r.width = width
	r.height = height
}

func (r *RightPanel) Focus() {}
func (r *RightPanel) Blur()  {}
