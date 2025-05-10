package ui

import "github.com/charmbracelet/lipgloss"

var Catppuccin = struct {
	Base      lipgloss.Color
	Surface0  lipgloss.Color
	Surface1  lipgloss.Color
	Surface2  lipgloss.Color
	Text      lipgloss.Color
	Subtext0  lipgloss.Color
	Subtext1  lipgloss.Color
	Overlay0  lipgloss.Color
	Overlay1  lipgloss.Color
	Overlay2  lipgloss.Color
	Blue      lipgloss.Color
	Green     lipgloss.Color
	Red       lipgloss.Color
	Yellow    lipgloss.Color
	Mauve     lipgloss.Color
	Pink      lipgloss.Color
	Flamingo  lipgloss.Color
	Rosewater lipgloss.Color
}{
	Base:      lipgloss.Color("#1e1e2e"),
	Surface0:  lipgloss.Color("#313244"),
	Surface1:  lipgloss.Color("#45475a"),
	Surface2:  lipgloss.Color("#585b70"),
	Text:      lipgloss.Color("#cdd6f4"),
	Subtext0:  lipgloss.Color("#a6adc8"),
	Subtext1:  lipgloss.Color("#bac2de"),
	Overlay0:  lipgloss.Color("#6c7086"),
	Overlay1:  lipgloss.Color("#7f849c"),
	Overlay2:  lipgloss.Color("#9399b2"),
	Blue:      lipgloss.Color("#89b4fa"),
	Green:     lipgloss.Color("#a6e3a1"),
	Red:       lipgloss.Color("#f38ba8"),
	Yellow:    lipgloss.Color("#f9e2af"),
	Mauve:     lipgloss.Color("#cba6f7"),
	Pink:      lipgloss.Color("#f5c2e7"),
	Flamingo:  lipgloss.Color("#f2cdcd"),
	Rosewater: lipgloss.Color("#f5e0dc"),
}

var (
	WindowStyle = lipgloss.NewStyle().
			Background(Catppuccin.Base).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(Catppuccin.Surface2).
			Padding(2, 4).
			Align(lipgloss.Center)

	TitleStyle = lipgloss.NewStyle().
			Foreground(Catppuccin.Blue).
			Bold(true).
			MarginBottom(1).
			Align(lipgloss.Center)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(Catppuccin.Base).
			Background(Catppuccin.Blue).
			Padding(0, 2).
			Bold(true).
			MarginTop(1).
			Align(lipgloss.Center)

	ButtonActiveStyle = ButtonStyle.Copy().
				Background(Catppuccin.Green)

	WordStyle = lipgloss.NewStyle().
			Foreground(Catppuccin.Text).
			Bold(true)

	TypedStyle = lipgloss.NewStyle().
			Foreground(Catppuccin.Overlay0)

	IncorrectStyle = lipgloss.NewStyle().
			Foreground(Catppuccin.Red)
)
