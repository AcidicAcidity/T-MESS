package tui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color
	Border     lipgloss.Color
	Error      lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
}

// Matrix — зелёный неон (классический хакерский стиль)
var Matrix = Theme{
	Name:       "matrix",
	Primary:    lipgloss.Color("#00ff00"),
	Secondary:  lipgloss.Color("#33ff33"),
	Accent:     lipgloss.Color("#66ff66"),
	Background: lipgloss.Color("#000000"),
	Foreground: lipgloss.Color("#00ff00"),
	Border:     lipgloss.Color("#00aa00"),
	Error:      lipgloss.Color("#ff0000"),
	Success:    lipgloss.Color("#00ff00"),
	Warning:    lipgloss.Color("#ffff00"),
}

// Cyberpunk — неон-розовый/голубой
var Cyberpunk = Theme{
	Name:       "cyberpunk",
	Primary:    lipgloss.Color("#00ffff"),
	Secondary:  lipgloss.Color("#ff00ff"),
	Accent:     lipgloss.Color("#ffff00"),
	Background: lipgloss.Color("#0a0a2a"),
	Foreground: lipgloss.Color("#00ffff"),
	Border:     lipgloss.Color("#ff00ff"),
	Error:      lipgloss.Color("#ff3366"),
	Success:    lipgloss.Color("#00ffaa"),
	Warning:    lipgloss.Color("#ffaa00"),
}

// Dark — спокойная тёмная
var Dark = Theme{
	Name:       "dark",
	Primary:    lipgloss.Color("#ffffff"),
	Secondary:  lipgloss.Color("#aaaaaa"),
	Accent:     lipgloss.Color("#3a6ea5"),
	Background: lipgloss.Color("#1e1e2e"),
	Foreground: lipgloss.Color("#cdd6f4"),
	Border:     lipgloss.Color("#313244"),
	Error:      lipgloss.Color("#f38ba8"),
	Success:    lipgloss.Color("#a6e3a1"),
	Warning:    lipgloss.Color("#f9e2af"),
}

// AvailableThemes — все доступные темы
var AvailableThemes = map[string]Theme{
	"matrix":    Matrix,
	"cyberpunk": Cyberpunk,
	"dark":      Dark,
}
