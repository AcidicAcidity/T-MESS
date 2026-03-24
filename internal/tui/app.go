package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AcidicAcidity/t-mess/internal/crypto"
	"github.com/AcidicAcidity/t-mess/internal/storage"
)

type App struct {
	splash *SplashScreen
	ready  bool
	width  int
	height int

	identity *crypto.Identity
	db       *storage.Database
	theme    Theme

	// Статус для отображения
	statusMsg  string
	statusTime time.Time
}

func NewApp(identity *crypto.Identity, db *storage.Database) *App {
	return &App{
		splash:     NewSplashScreen(),
		identity:   identity,
		db:         db,
		theme:      Matrix,
		statusMsg:  "Initializing...",
		statusTime: time.Now(),
	}
}

func (a *App) Init() tea.Cmd {
	return a.splash.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.splash.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		if !a.ready {
			return a, nil
		}
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "s":
			a.statusMsg = fmt.Sprintf("Node: %s", a.identity.PeerID[:16]+"...")
			a.statusTime = time.Now()
		}
	}

	if !a.ready {
		var cmd tea.Cmd
		a.splash, cmd = a.splash.Update(msg)
		if a.splash.IsDone() {
			a.ready = true
			a.statusMsg = fmt.Sprintf("Node: %s", a.identity.PeerID[:16]+"...")
		}
		return a, cmd
	}

	return a, nil
}

func (a *App) View() string {
	if !a.ready {
		return a.splash.View()
	}

	// Статусная строка
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		MarginBottom(1)

	// Основное окно
	infoStyle := lipgloss.NewStyle().
		Foreground(a.theme.Primary).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(a.theme.Border).
		Padding(1, 2)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		statusStyle.Render(fmt.Sprintf("> %s", a.statusMsg)),
		infoStyle.Render("Welcome to T-MESS!"),
		"",
		lipgloss.NewStyle().Foreground(a.theme.Secondary).Render("Node ID:"),
		a.identity.PeerID,
		"",
		lipgloss.NewStyle().Foreground(a.theme.Secondary).Render("Fingerprint:"),
		a.identity.Fingerprint(),
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Press Ctrl+C to quit | s - show node info"),
	)

	return lipgloss.Place(
		a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (a *App) Run() error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}
