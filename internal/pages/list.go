package pages

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/luynrs/justssh/internal/models"
	"github.com/luynrs/justssh/internal/theme"
)

type ListModel struct {
	styles    theme.Styles
	servers   []models.Server
	cursor    int
	searching bool
	search    textinput.Model
	termWidth int

	confirm bool
}

func NewList(styles theme.Styles, servers []models.Server) ListModel {
	search := textinput.New()
	search.Prompt = ""
	search.Placeholder = "type to filter..."

	return ListModel{styles: styles, servers: servers, search: search}
}

func (m *ListModel) SetServers(servers []models.Server) {
	m.servers = servers
	if max := len(m.visible()) - 1; m.cursor > max {
		m.cursor = max
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m ListModel) IsSearching() bool {
	return m.searching
}

func (m *ListModel) SetTermWidth(w int) {
	m.termWidth = w
}

func (m ListModel) Update(msg tea.KeyMsg) (ListModel, tea.Cmd) {
	if m.confirm {
		m.confirm = false
		if msg.String() == "y" {
			return m, m.deleteCmd()
		}
		return m, nil
	}
	if m.searching {
		return m.updateSearching(msg)
	}

	switch msg.String() {
	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)
	case "enter":
		return m, m.connectCmd()
	case "/":
		m.searching = true
		return m, m.search.Focus()
	case "a":
		return m, func() tea.Msg { return AddRequested{} }
	case "e":
		return m, m.editCmd()
	case "d":
		if _, ok := m.selected(); ok {
			m.confirm = true
		}
	}
	return m, nil
}

func (m ListModel) updateSearching(msg tea.KeyMsg) (ListModel, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.searching = false
		m.search.SetValue("")
		m.search.Blur()
		m.cursor = 0
		return m, nil
	case tea.KeyEnter:
		return m, m.connectCmd()
	case tea.KeyUp:
		m.moveCursor(-1)
		return m, nil
	case tea.KeyDown:
		m.moveCursor(1)
		return m, nil
	default:
		var cmd tea.Cmd
		m.search, cmd = m.search.Update(msg)
		m.cursor = 0
		return m, cmd
	}
}

func (m *ListModel) moveCursor(delta int) {
	n := len(m.visible())
	if n == 0 {
		m.cursor = 0
		return
	}
	m.cursor = (m.cursor + delta + n) % n
}

func (m ListModel) visible() []int {
	query := strings.ToLower(strings.TrimSpace(m.search.Value()))
	idx := make([]int, 0, len(m.servers))
	for i, s := range m.servers {
		if query == "" || strings.Contains(strings.ToLower(s.Name), query) ||
			strings.Contains(strings.ToLower(s.Host), query) ||
			strings.Contains(strings.ToLower(s.User), query) {
			idx = append(idx, i)
		}
	}
	return idx
}

func (m ListModel) selected() (int, bool) {
	idx := m.visible()
	if len(idx) == 0 {
		return 0, false
	}
	return idx[m.cursor], true
}

func (m ListModel) connectCmd() tea.Cmd {
	i, ok := m.selected()
	if !ok {
		return nil
	}
	server := m.servers[i]
	return func() tea.Msg { return ConnectRequested{Server: server} }
}

func (m ListModel) editCmd() tea.Cmd {
	i, ok := m.selected()
	if !ok {
		return nil
	}
	server := m.servers[i]
	return func() tea.Msg { return EditRequested{Index: i, Server: server} }
}

func (m ListModel) deleteCmd() tea.Cmd {
	i, ok := m.selected()
	if !ok {
		return nil
	}
	return func() tea.Msg { return DeleteRequested{Index: i} }
}

func (m ListModel) View() string {
	var b strings.Builder

	idx := m.visible()

	footer := m.footer()
	alignWidth := m.alignWidth(footer)
	if m.confirm && len(idx) > 0 {
		footer = m.confirmFooter(m.servers[idx[m.cursor]].Name, alignWidth)
	}
	b.WriteString(m.header(len(idx), alignWidth))
	b.WriteString("\n\n")

	if len(idx) == 0 {
		b.WriteString(m.styles.Dim.Render("  No servers. Press a to add one."))
		b.WriteString("\n")
	}

	width := nameWidth(m.servers)
	for i, si := range idx {
		s := m.servers[si]
		name := padRight(s.Name, width)
		cursor := "  "
		rendered := m.styles.Name.Render(name)
		if i == m.cursor {
			cursor = m.styles.Cursor.Render("❯ ")
			rendered = m.styles.Selected.Render(name)
		}
		row := alignRight(cursor+rendered, m.styles.Dim.Render(target(s)), alignWidth)
		b.WriteString(row + "\n")
	}

	b.WriteString("\n")
	b.WriteString(footer)

	return b.String()
}

// natural width from the footer text, clamped to the terminal so long
// lines shrink instead of overflowing it
func (m ListModel) alignWidth(footer string) int {
	w := lipgloss.Width(footer)
	if m.termWidth > 0 && m.termWidth-1 < w {
		w = m.termWidth - 1
	}
	return w
}

func (m ListModel) header(visibleCount, footerWidth int) string {
	title := m.styles.Title.Render("JustSSH")
	count := m.styles.Dim.Render(fmt.Sprintf("%d/%d", visibleCount, len(m.servers)))

	if m.searching {
		left := title + " " + m.styles.Dim.Render("~ Search:") + " " + m.search.View()
		return alignRight(left, count, footerWidth)
	}

	return alignRight(title, count, footerWidth)
}

func alignRight(left, right string, width int) string {
	pad := width - lipgloss.Width(left) - lipgloss.Width(right)
	if pad < 2 {
		pad = 2
	}
	return left + strings.Repeat(" ", pad) + right
}

func (m ListModel) footer() string {
	return renderHints(m.styles, [][2]string{
		{"↵", "Connect"},
		{"/", "Search"},
		{"a", "Add"},
		{"e", "Edit"},
		{"d", "Delete"},
		{"q", "Quit"},
	})
}

func (m ListModel) confirmFooter(name string, width int) string {
	msg := m.styles.Error.Render(fmt.Sprintf("Delete %s?", name))
	hints := renderHints(m.styles, [][2]string{{"y", "Confirm"}, {"any", "Cancel"}})
	return alignRight(msg, hints, width)
}

func renderHints(styles theme.Styles, hints [][2]string) string {
	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = styles.Key.Render(h[0]) + " " + styles.Help.Render(h[1])
	}
	return strings.Join(parts, "    ")
}

func target(s models.Server) string {
	if s.User == "" {
		return s.Host
	}
	return s.User + "@" + s.Host
}

func nameWidth(servers []models.Server) int {
	width := 0
	for _, s := range servers {
		if n := utf8.RuneCountInString(s.Name); n > width {
			width = n
		}
	}
	return width
}

func padRight(s string, width int) string {
	if pad := width - utf8.RuneCountInString(s); pad > 0 {
		return s + strings.Repeat(" ", pad)
	}
	return s
}
