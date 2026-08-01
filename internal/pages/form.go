package pages

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/luynrs/justssh/internal/models"
	"github.com/luynrs/justssh/internal/theme"
)

type field int

const (
	fieldName field = iota
	fieldHost
	fieldUser
	fieldPort
	fieldKey
	fieldCount
)

var fieldLabels = [fieldCount]string{
	fieldName: "Name",
	fieldHost: "Host",
	fieldUser: "User",
	fieldPort: "Port",
	fieldKey:  "Key",
}

type FormModel struct {
	styles    theme.Styles
	title     string
	inputs    [fieldCount]textinput.Model
	focus     field
	editIndex *int // nil: append, non-nil: replace at index
	err       string
}

func NewForm(styles theme.Styles) FormModel {
	return newForm(styles, "Add server", nil, models.Server{Port: 22})
}

func NewEditForm(styles theme.Styles, index int, server models.Server) FormModel {
	return newForm(styles, "Edit server", &index, server)
}

func newForm(styles theme.Styles, title string, editIndex *int, s models.Server) FormModel {
	m := FormModel{styles: styles, title: title, editIndex: editIndex}

	values := [fieldCount]string{
		fieldName: s.Name,
		fieldHost: s.Host,
		fieldUser: s.User,
		fieldPort: portString(s.Port),
		fieldKey:  s.Key,
	}
	for f := field(0); f < fieldCount; f++ {
		in := textinput.New()
		in.Prompt = ""
		in.SetValue(values[f])
		m.inputs[f] = in
	}
	m.inputs[fieldName].Focus()
	return m
}

func (m FormModel) EditIndex() *int {
	return m.editIndex
}

func portString(port int) string {
	if port == 0 {
		return ""
	}
	return strconv.Itoa(port)
}

func (m FormModel) Update(msg tea.KeyMsg) (FormModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return m, func() tea.Msg { return FormCancelled{} }
	case "enter":
		return m.submit()
	case "tab", "down":
		m.setFocus(m.focus + 1)
		return m, nil
	case "shift+tab", "up":
		m.setFocus(m.focus - 1 + fieldCount)
		return m, nil
	}

	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	return m, cmd
}

func (m *FormModel) setFocus(f field) {
	m.inputs[m.focus].Blur()
	m.focus = f % fieldCount
	m.inputs[m.focus].Focus()
}

func (m FormModel) submit() (FormModel, tea.Cmd) {
	name := strings.TrimSpace(m.inputs[fieldName].Value())
	host := strings.TrimSpace(m.inputs[fieldHost].Value())
	if name == "" || host == "" {
		m.err = "Name and Host are required"
		return m, nil
	}

	port := 22
	if raw := strings.TrimSpace(m.inputs[fieldPort].Value()); raw != "" {
		p, err := strconv.Atoi(raw)
		if err != nil || p <= 0 {
			m.err = "Port must be a positive number"
			return m, nil
		}
		port = p
	}

	server := models.Server{
		Name: name,
		Host: host,
		User: strings.TrimSpace(m.inputs[fieldUser].Value()),
		Port: port,
		Key:  strings.TrimSpace(m.inputs[fieldKey].Value()),
	}
	return m, func() tea.Msg { return FormSaved{Server: server} }
}

func (m FormModel) View() string {
	var b strings.Builder

	b.WriteString(m.styles.Title.Render(m.title))
	b.WriteString("\n\n")

	for f := field(0); f < fieldCount; f++ {
		b.WriteString(m.styles.Label.Render(fieldLabels[f]))
		b.WriteString(m.inputs[f].View())
		b.WriteString("\n")
	}

	if m.err != "" {
		b.WriteString("\n")
		b.WriteString(m.styles.Error.Render(m.err))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(renderHints(m.styles, [][2]string{
		{"↵", "Save"},
		{"Esc", "Cancel"},
	}))

	return b.String()
}
