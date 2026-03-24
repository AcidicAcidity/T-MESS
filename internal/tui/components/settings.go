package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Settings struct {
	selected  int
	nick      string
	themeName string
	language  string
	width     int
	height    int
	theme     lipgloss.Color
	onSave    func(nick, theme, lang string)
	onBack    func()
}

type SettingsItem struct {
	Name  string
	Value string
}

var themes = []string{"matrix", "cyberpunk", "dark"}
var languages = []string{"English", "Русский"}

func NewSettings(theme lipgloss.Color, currentNick, currentTheme, currentLang string) *Settings {
	return &Settings{
		selected:  0,
		nick:      currentNick,
		themeName: currentTheme,
		language:  currentLang,
		theme:     theme,
	}
}

func (s *Settings) Update(msg tea.Msg) (*Settings, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			s.selected--
			if s.selected < 0 {
				s.selected = 3
			}
			return s, nil
		case "down", "j":
			s.selected++
			if s.selected > 3 {
				s.selected = 0
			}
			return s, nil
		case "left", "h":
			s.changeValue(-1)
			return s, nil
		case "right", "l":
			s.changeValue(1)
			return s, nil
		case "enter":
			if s.selected == 3 { // Save & Back
				if s.onSave != nil {
					s.onSave(s.nick, s.themeName, s.language)
				}
				if s.onBack != nil {
					s.onBack()
				}
			}
			return s, nil
		case "esc":
			if s.onBack != nil {
				s.onBack()
			}
			return s, nil
		}
	}
	return s, nil
}

func (s *Settings) changeValue(delta int) {
	switch s.selected {
	case 0: // Nickname
		// Для никнейма используем отдельный ввод
	case 1: // Theme
		idx := 0
		for i, t := range themes {
			if t == s.themeName {
				idx = i
				break
			}
		}
		idx += delta
		if idx < 0 {
			idx = len(themes) - 1
		}
		if idx >= len(themes) {
			idx = 0
		}
		s.themeName = themes[idx]
	case 2: // Language
		idx := 0
		for i, l := range languages {
			if l == s.language {
				idx = i
				break
			}
		}
		idx += delta
		if idx < 0 {
			idx = len(languages) - 1
		}
		if idx >= len(languages) {
			idx = 0
		}
		s.language = languages[idx]
	}
}

func (s *Settings) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(s.theme).
		Bold(true).
		MarginBottom(2)

	items := []SettingsItem{
		{Name: "Nickname", Value: s.nick},
		{Name: "Theme", Value: s.themeName},
		{Name: "Language", Value: s.language},
		{Name: "Save & Back", Value: ""},
	}

	var lines []string
	for i, item := range items {
		var line string
		if i == s.selected {
			if i == 3 {
				line = lipgloss.NewStyle().
					Foreground(s.theme).
					Bold(true).
					Render(fmt.Sprintf("▶  %s", item.Name))
			} else {
				line = lipgloss.NewStyle().
					Foreground(s.theme).
					Bold(true).
					Render(fmt.Sprintf("▶  %s: %s  ◀", item.Name, item.Value))
			}
		} else {
			if i == 3 {
				line = fmt.Sprintf("   %s", item.Name)
			} else {
				line = fmt.Sprintf("   %s: %s", item.Name, item.Value)
			}
		}
		lines = append(lines, line)
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		MarginTop(2)

	help := helpStyle.Render("↑/k  ↓/j  navigate • ←/h  →/l  change • Enter save • Esc back")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("⚙️  SETTINGS"),
		lipgloss.JoinVertical(lipgloss.Left, lines...),
		help,
	)

	return lipgloss.Place(
		s.width, s.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (s *Settings) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *Settings) SetOnSave(f func(string, string, string)) {
	s.onSave = f
}

func (s *Settings) SetOnBack(f func()) {
	s.onBack = f
}
