package theme

import "github.com/charmbracelet/lipgloss"

// ansi 0-15
const (
	accent = lipgloss.Color("6") // cyan
	dim    = lipgloss.Color("8") // bright black / gray
	red    = lipgloss.Color("1")
)

type Styles struct {
	Title    lipgloss.Style
	Cursor   lipgloss.Style
	Name     lipgloss.Style
	Selected lipgloss.Style
	Dim      lipgloss.Style
	Key      lipgloss.Style
	Help     lipgloss.Style
	Error    lipgloss.Style
	Label    lipgloss.Style
}

func New() Styles {
	return Styles{
		Title:    lipgloss.NewStyle().Bold(true).Foreground(accent),
		Cursor:   lipgloss.NewStyle().Foreground(accent).Bold(true),
		Name:     lipgloss.NewStyle(),
		Selected: lipgloss.NewStyle().Bold(true).Foreground(accent),
		Dim:      lipgloss.NewStyle().Foreground(dim),
		Key:      lipgloss.NewStyle().Bold(true).Foreground(accent),
		Help:     lipgloss.NewStyle().Foreground(dim),
		Error:    lipgloss.NewStyle().Foreground(red),
		Label:    lipgloss.NewStyle().Foreground(dim).Width(9),
	}
}
