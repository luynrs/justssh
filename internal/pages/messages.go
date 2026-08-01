package pages

import "github.com/luynrs/justssh/internal/models"

type (
	ConnectRequested struct{ Server models.Server }
	AddRequested     struct{}
)

type EditRequested struct {
	Index  int
	Server models.Server
}

type (
	DeleteRequested struct{ Index int }
	FormSaved       struct{ Server models.Server }
	FormCancelled   struct{}
)
