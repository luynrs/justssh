package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/luynrs/justssh/internal/models"
	"github.com/luynrs/justssh/internal/pages"
	"github.com/luynrs/justssh/internal/storage"
	"github.com/luynrs/justssh/internal/theme"
)

type page int

const (
	pageList page = iota
	pageAdd
	pageEdit
)

type Model struct {
	styles  theme.Styles
	store   *storage.Store
	servers []models.Server

	page page
	list pages.ListModel
	form pages.FormModel

	err      error
	quitting bool
}

func New(store *storage.Store) (Model, error) {
	servers, err := store.Load()
	if err != nil {
		return Model{}, err
	}

	styles := theme.New()
	return Model{
		styles:  styles,
		store:   store,
		servers: servers,
		list:    pages.NewList(styles, servers),
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}
