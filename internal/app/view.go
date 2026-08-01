package app

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	out := ""
	switch m.page {
	case pageList:
		out = m.list.View()
	case pageAdd, pageEdit:
		out = m.form.View()
	}

	if m.err != nil {
		out += "\n\n" + m.styles.Error.Render(m.err.Error())
	}

	return out
}
