package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SplashScreen struct {
	phase    int
	progress int
	messages []string
	done     bool
}

var asciiLogo = `
╔═══════════════════════════════════════════════════════════════════╗
║                                                                   ║
║     ████████╗ ██████╗ ███╗   ███╗███████╗███████╗███████╗        ║
║     ╚══██╔══╝██╔═══██╗████╗ ████║██╔════╝██╔════╝██╔════╝        ║
║        ██║   ██║   ██║██╔████╔██║█████╗  ███████╗███████╗        ║
║        ██║   ██║   ██║██║╚██╔╝██║██╔══╝  ╚════██║╚════██║        ║
║        ██║   ╚██████╔╝██║ ╚═╝ ██║███████╗███████║███████║        ║
║        ╚═╝    ╚═════╝ ╚═╝     ╚═╝╚══════╝╚══════╝╚══════╝        ║
║                                                                   ║
║              [ TERMINAL MESSENGER v0.1.0 ]                       ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
`

func NewSplashScreen() *SplashScreen {
	return &SplashScreen{
		messages: []string{
			"INITIALIZING CRYPTO ENGINE.........",
			"GENERATING NODE IDENTITY............",
			"CONNECTING TO DHT...................",
			"BOOTSTRAPPING NETWORK...............",
			"SYNCING DISTRIBUTED HISTORY.........",
			"LOADING INTERFACE...................",
		},
	}
}

func (s *SplashScreen) Init() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return splashTickMsg{}
	})
}

type splashTickMsg struct{}

func (s *SplashScreen) Update(msg tea.Msg) (*SplashScreen, tea.Cmd) {
	switch msg.(type) {
	case splashTickMsg:
		if s.phase < len(s.messages) {
			s.phase++
			return s, tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
				return splashTickMsg{}
			})
		} else if s.progress < 100 {
			s.progress += 10
			if s.progress >= 100 {
				s.done = true
				return s, nil
			}
			return s, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
				return splashTickMsg{}
			})
		}
	}
	return s, nil
}

func (s *SplashScreen) View() string {
	if s.done {
		return ""
	}

	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00")).
		Bold(true)

	// Список инициализации
	var initLog strings.Builder
	for i := 0; i < s.phase && i < len(s.messages); i++ {
		status := "✓"
		if i == s.phase-1 {
			status = ">"
		}
		initLog.WriteString(fmt.Sprintf("  %s %s\n", status, s.messages[i]))
	}

	// Прогресс-бар
	barWidth := 50
	filled := int(float64(barWidth) * float64(s.progress) / 100)
	progressBar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00"))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logoStyle.Render(asciiLogo),
		"",
		initLog.String(),
		"",
		progressStyle.Render(progressBar),
		fmt.Sprintf("  %d%%", s.progress),
	)

	return lipgloss.Place(
		80, 30,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}
