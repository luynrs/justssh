package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/luynrs/justssh/internal/storage"
)

func Run(store *storage.Store) error {
	model, err := New(store)
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(model).Run()
	return err
}
