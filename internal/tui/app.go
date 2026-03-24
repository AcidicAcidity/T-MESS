package tui

import (
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
}

func NewApp(identity *crypto.Identity, db *storage.Database) *App {
	return &App{
		splash:   NewSplashScreen(),
		identity: identity,
		db:       db,
		theme:    Matrix, // по умолчанию матричная тема
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

	case tea.KeyMsg:
		if !a.ready && msg.String() == "enter" {
			a.ready = true
			return a, nil
		}
		if a.ready && msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
	}

	if !a.ready {
		var cmd tea.Cmd
		a.splash, cmd = a.splash.Update(msg)
		return a, cmd
	}

	return a, nil
}

func (a *App) View() string {
	if !a.ready {
		return a.splash.View()
	}

	// Временное главное окно, пока не добавили компоненты
	infoStyle := lipgloss.NewStyle().
		Foreground(a.theme.Primary).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(a.theme.Border).
		Padding(1, 2)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		infoStyle.Render("Welcome to T-MESS!"),
		"",
		lipgloss.NewStyle().Foreground(a.theme.Secondary).Render("Node ID:"),
		a.identity.PeerID,
		"",
		lipgloss.NewStyle().Foreground(a.theme.Secondary).Render("Fingerprint:"),
		a.identity.Fingerprint(),
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Press Ctrl+C to quit"),
	)

	return lipgloss.Place(
		a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// Run запускает TUI приложение
func (a *App) Run() error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}
