package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuChoice int

const (
	MenuConnectLocal MenuChoice = iota
	MenuConnectGlobal
	MenuSettings
	MenuExit
)

type Menu struct {
	selected int
	width    int
	height   int
	theme    lipgloss.Color
	nick     string
}

type MenuItem struct {
	Name string
	Icon string
	Desc string
}

var menuItems = []MenuItem{
	{Name: "Connect Local", Icon: "🌐", Desc: "Find and connect to peers in your local network"},
	{Name: "Connect Global", Icon: "🌍", Desc: "Connect to the global P2P network"},
	{Name: "Settings", Icon: "⚙️", Desc: "Change theme, language, nickname"},
	{Name: "Exit", Icon: "🚪", Desc: "Exit T-MESS"},
}

func NewMenu(theme lipgloss.Color, nick string) *Menu {
	return &Menu{
		selected: 0,
		theme:    theme,
		nick:     nick,
	}
}

func (m *Menu) Update(msg tea.Msg) (*Menu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.selected--
			if m.selected < 0 {
				m.selected = len(menuItems) - 1
			}
			return m, nil
		case "down", "j":
			m.selected++
			if m.selected >= len(menuItems) {
				m.selected = 0
			}
			return m, nil
		case "enter", " ":
			return m, m.selectItem()
		}
	}
	return m, nil
}

func (m *Menu) selectItem() tea.Cmd {
	switch m.selected {
	case int(MenuConnectLocal):
		return func() tea.Msg { return MenuSelectedMsg{Choice: MenuConnectLocal} }
	case int(MenuConnectGlobal):
		return func() tea.Msg { return MenuSelectedMsg{Choice: MenuConnectGlobal} }
	case int(MenuSettings):
		return func() tea.Msg { return MenuSelectedMsg{Choice: MenuSettings} }
	case int(MenuExit):
		return func() tea.Msg { return ExitRequestMsg{} }
	}
	return nil
}

func (m *Menu) View() string {
	// Заголовок с ASCII артом
	logoStyle := lipgloss.NewStyle().
		Foreground(m.theme).
		Bold(true).
		MarginBottom(1)

	logo := `
╔═══════════════════════════════════════════════════════════════════╗
║                                                                   ║
║     ████████╗ ██████╗ ███╗   ███╗███████╗███████╗███████╗        ║
║     ╚══██╔══╝██╔═══██╗████╗ ████║██╔════╝██╔════╝██╔════╝        ║
║        ██║   ██║   ██║██╔████╔██║█████╗  ███████╗███████╗        ║
║        ██║   ██║   ██║██║╚██╔╝██║██╔══╝  ╚════██║╚════██║        ║
║        ██║   ╚██████╔╝██║ ╚═╝ ██║███████╗███████║███████║        ║
║        ╚═╝    ╚═════╝ ╚═╝     ╚═╝╚══════╝╚══════╝╚══════╝        ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝`

	// Приветствие
	greeting := fmt.Sprintf("Welcome back, %s!", m.nick)
	greetingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00")).
		Italic(true).
		MarginBottom(2)

	// Меню
	var menuLines []string
	for i, item := range menuItems {
		// Иконка и имя
		icon := item.Icon
		name := item.Name

		// Стиль для выбранного пункта
		var line string
		if i == m.selected {
			line = lipgloss.NewStyle().
				Foreground(m.theme).
				Bold(true).
				PaddingLeft(4).
				Render(fmt.Sprintf("▶  %s %s", icon, name))
		} else {
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				PaddingLeft(6).
				Render(fmt.Sprintf("  %s %s", icon, name))
		}
		menuLines = append(menuLines, line)
	}

	menuStyle := lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(2)

	menuContent := menuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, menuLines...))

	// Подсказки
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		MarginTop(2)

	help := helpStyle.Render("↑/k  ↓/j  to navigate • Enter to select • Ctrl+C to exit")

	// Сборка
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logoStyle.Render(logo),
		greetingStyle.Render(greeting),
		menuContent,
		help,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m *Menu) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Menu) SetNick(nick string) {
	m.nick = nick
}

type MenuSelectedMsg struct {
	Choice MenuChoice
}

type ExitRequestMsg struct{}
