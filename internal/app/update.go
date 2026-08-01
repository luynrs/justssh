package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/luynrs/justssh/internal/models"
	"github.com/luynrs/justssh/internal/pages"
	"github.com/luynrs/justssh/internal/ssh"
)

type connectFinishedMsg struct{ err error }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.updateKey(msg)

	case tea.WindowSizeMsg:
		m.list.SetTermWidth(msg.Width)
		return m, nil

	case pages.ConnectRequested:
		return m, m.connectCmd(msg.Server)

	case connectFinishedMsg:
		m.err = msg.err
		return m, nil

	case pages.AddRequested:
		m.page = pageAdd
		m.form = pages.NewForm(m.styles)
		return m, nil

	case pages.EditRequested:
		m.page = pageEdit
		m.form = pages.NewEditForm(m.styles, msg.Index, msg.Server)
		return m, nil

	case pages.DeleteRequested:
		m.servers = append(m.servers[:msg.Index:msg.Index], m.servers[msg.Index+1:]...)
		m.persist()
		return m, nil

	case pages.FormSaved:
		if idx := m.form.EditIndex(); idx != nil {
			m.servers[*idx] = msg.Server
		} else {
			m.servers = append(m.servers, msg.Server)
		}
		m.persist()
		m.page = pageList
		return m, nil

	case pages.FormCancelled:
		m.page = pageList
		return m, nil
	}

	return m, nil
}

// fast typing split before disp.
func (m Model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyRunes && len(msg.Runes) > 1 && !msg.Paste {
		var cmds []tea.Cmd
		for _, r := range msg.Runes {
			var cmd tea.Cmd
			var next tea.Model
			next, cmd = m.updateKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}, Alt: msg.Alt})
			m = next.(Model)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)
	}

	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "q":
		if m.page == pageList && !m.list.IsSearching() {
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.page {
	case pageList:
		m.list, cmd = m.list.Update(msg)
	case pageAdd, pageEdit:
		m.form, cmd = m.form.Update(msg)
	}
	return m, cmd
}

func (m *Model) persist() {
	m.err = m.store.Save(m.servers)
	m.list.SetServers(m.servers)
}

func (m Model) connectCmd(s models.Server) tea.Cmd {
	return tea.ExecProcess(ssh.Command(s), func(err error) tea.Msg {
		return connectFinishedMsg{err: err}
	})
}
