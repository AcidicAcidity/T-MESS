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
	userNick         string
	peerID           string
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

func (r *RightPanel) SetUserNick(nick string) {
	r.userNick = nick
}

func (r *RightPanel) SetPeerID(peerID string) {
	r.peerID = peerID
}

func (r *RightPanel) Update(msg tea.Msg) (*RightPanel, tea.Cmd) {
	return r, nil
}

func (r *RightPanel) View() string {
	if r.chat == nil {
		return ""
	}

	peerDisplay := ""
	if r.peerID != "" {
		peerShort := r.peerID
		if len(peerShort) > 16 {
			peerShort = peerShort[:14] + ".."
		}
		peerDisplay = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render(peerShort)
	}

	// Инфо о текущем пользователе
	userSection := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("👤 YOUR IDENTITY"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Render(r.userNick),
		peerDisplay,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render(r.peerID[:16]+"..."),
		"",
	)

	// Инфо о чате
	chatSection := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("💬 CURRENT CHAT"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Render(r.chat.Name),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render(fmt.Sprintf("Type: %s", r.chat.Type)),
		"",
	)

	// Статус подключения
	statusColor := lipgloss.Color("#ff0000")
	if strings.Contains(r.connectionStatus, "ONLINE") {
		statusColor = lipgloss.Color("#00ff00")
	} else if strings.Contains(r.connectionStatus, "LOCAL") {
		statusColor = lipgloss.Color("#ffff00")
	}

	statusSection := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("🔌 CONNECTION"),
		lipgloss.NewStyle().
			Foreground(statusColor).
			Render(r.connectionStatus),
		"",
	)

	// Эмодзи (заглушка)
	emojiSection := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("😊 QUICK EMOJIS"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffff00")).
			Render("😊 😂 🔥 🚀 💯 ❤️ 🎉 🤔"),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		userSection,
		chatSection,
		statusSection,
		emojiSection,
	)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Width(r.width).
		Render(content)
}

func (r *RightPanel) SetSize(width, height int) {
	r.width = width
	r.height = height
}

func (r *RightPanel) Focus() {}
func (r *RightPanel) Blur()  {}
