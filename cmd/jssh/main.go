package main

import (
	"fmt"
	"os"

	"github.com/luynrs/justssh/internal/app"
	"github.com/luynrs/justssh/internal/storage"
)

func main() {
	store, err := storage.NewStore()
	if err != nil {
		fmt.Fprintln(os.Stderr, "jssh:", err)
		os.Exit(1)
	}

	if err := app.Run(store); err != nil {
		fmt.Fprintln(os.Stderr, "jssh:", err)
		os.Exit(1)
	}
}
